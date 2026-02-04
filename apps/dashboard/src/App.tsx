import { createEffect, createSignal, onCleanup, onMount, createMemo, untrack, Show } from "solid-js"
import Header from "./components/Header"
import Sidebar from "./components/Sidebar"
import TaskTable from "./components/TaskTable"
import Editor from "./components/Editor"
import "./App.css"
import { keyArray } from "@solid-primitives/keyed"
import { ReactiveMap } from "@solid-primitives/map"
import { sortedIndex, sortedIndexCmp, sortedInsert, sortedRemove } from "./lib"
import { createMutable, createStore } from "solid-js/store"
import { MyTransitionGroup } from "./components/MyTransitionGroup"
import "./components/TaskTable.css"

type Tab = "tasks" | "roles" | "templates"

export type TaskItem = {
  id: string
  short_id: string
  title: string
  role: string
  priority: string
  completed: boolean
  parent: string
  blockers: string[]
  blocks: string[]
  path: string
  date_created: string
  date_edited: string
}

export type TaskTreeNode = {
  task: TaskItem
  depth: number
  expanded?: boolean
  lastUpdateTick: number
}

export type SortField = "title" | "priority" | "date_created" | "date_edited" | "short_id" // | "actionable"
export type SortDirection = "asc" | "desc"

type FileEntry = {
  name: string
  path: string
  kind: string
}

type FilePayload = {
  path: string
  content: string
}

type StreamUpdate = {
  event: string
  path: string
  project: string
  task?: {
    id: string
    file_path: string
    title: string
  }
}

type ProjectInfo = {
  name: string
  storage_root: string
  tasks_root: string
  roles_root: string
  templates_root: string
  git_root: string
  storage: string
}

type ProjectsResponse = {
  projects: ProjectInfo[]
  current: string
}

type TaskFrontmatter = {
  title?: string
  role?: string
  priority?: string
  status?: string
  completed?: boolean
  blockers?: string[]
  blocks?: string[]
  parent?: string
}

type ParsedTask = {
  frontmatter: TaskFrontmatter
  body: string
}

function parseFrontmatter(content: string): ParsedTask {
  const lines = content.split("\n")
  if (lines[0] !== "---") {
    return { frontmatter: {}, body: content }
  }

  let endIndex = -1
  for (let i = 1; i < lines.length; i++) {
    if (lines[i] === "---") {
      endIndex = i
      break
    }
  }

  if (endIndex === -1) {
    return { frontmatter: {}, body: content }
  }

  const frontmatterLines = lines.slice(1, endIndex)
  const bodyLines = lines.slice(endIndex + 1).join("\n")
  const frontmatter: TaskFrontmatter = {}

  for (const line of frontmatterLines) {
    const colonIndex = line.indexOf(":")
    if (colonIndex === -1) continue
    const key = line.substring(0, colonIndex).trim()
    const value = line.substring(colonIndex + 1).trim()

    if (value === "true") {
      (frontmatter as any)[key] = true
    } else if (value === "false") {
      (frontmatter as any)[key] = false
    } else if (value.startsWith("[") && value.endsWith("]")) {
      try {
        (frontmatter as any)[key] = JSON.parse(value)
      } catch {
        (frontmatter as any)[key] = value
      }
    } else {
      (frontmatter as any)[key] = value
    }
  }

  return { frontmatter, body: bodyLines.trim() }
}

async function fetchJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Request failed: ${res.status}`)
  }
  return res.json() as Promise<T>
}

function errorMessage(err: unknown) {
  if (err instanceof Error) {
    return err.message
  }
  return "Unknown error"
}

export default function App() {
  const [tab, setTab] = createSignal<Tab>("tasks")
  const [tasks, setTasks] = createSignal<TaskItem[]>([])
  const [activePath, setActivePath] = createSignal("")
  const [content, setContent] = createSignal("")
  const [dirty, setDirty] = createSignal(false)
  const [status, setStatus] = createSignal("")
  const [connected, setConnected] = createSignal(false)
  const [lastEvent, setLastEvent] = createSignal("")
  const [projects, setProjects] = createSignal<ProjectInfo[]>([])
  const [currentProject, setCurrentProject] = createSignal("")

  // Task filtering, sorting, and search
  const [searchQuery, setSearchQuery] = createSignal("")
  const [filterStatus, setFilterStatus] = createSignal<"all" | "active" | "done">("done")
  const [filterRole, setFilterRole] = createSignal<string>("all")
  const [filterPriority, setFilterPriority] = createSignal<string>("all")
  const [sortField, setSortField] = createSignal<SortField>("priority")
  const [sortDirection, setSortDirection] = createSignal<SortDirection>("desc")
  const [parsedTask, setParsedTask] = createSignal<ParsedTask>({ frontmatter: {}, body: "" })

  const apiURL = (path: string) => {
    const project = currentProject()
    if (!project) return path
    const separator = path.includes('?') ? '&' : '?'
    return `${path}${separator}project=${encodeURIComponent(project)}`
  }

  const loadProjects = async () => {
    try {
      const data = await fetchJSON<ProjectsResponse>("/api/projects")
      setProjects(data.projects)
      if (!currentProject() && data.current) {
        setCurrentProject(data.current)
      }
      const params = new URLSearchParams(window.location.search)
      const urlProject = params.get('project')
      if (urlProject) {
        setCurrentProject(urlProject)
      }
    } catch (err) {
      setStatus(`Failed to load projects: ${errorMessage(err)}`)
    }
  }

  const loadTasks = async () => {
    try {
      const data = await fetchJSON<TaskItem[]>(apiURL("/api/tasks"))
      setTasks(data)
    } catch (err) {
      setStatus(`Failed to load tasks: ${errorMessage(err)}`)
    }
  }

  const loadFile = async (path: string) => {
    try {
      const data = await fetchJSON<FilePayload>(apiURL(`/api/file?path=${encodeURIComponent(path)}`))
      setActivePath(data.path)
      setContent(data.content)
      setParsedTask(parseFrontmatter(data.content))
      setDirty(false)
      setStatus(`Loaded ${data.path}`)
    } catch (err) {
      setStatus(`Failed to load file: ${errorMessage(err)}`)
    }
  }

  const updateContent = () => {
    const parsed = parsedTask()
    let result = "---\n"

    if (parsed.frontmatter.title) result += `title: ${parsed.frontmatter.title}\n`
    if (parsed.frontmatter.role) result += `role: ${parsed.frontmatter.role}\n`
    if (parsed.frontmatter.priority) result += `priority: ${parsed.frontmatter.priority}\n`
    if (parsed.frontmatter.completed !== undefined) {
      result += `completed: ${parsed.frontmatter.completed ? "true" : "false"}\n`
    }
    if (parsed.frontmatter.blockers && parsed.frontmatter.blockers.length > 0) {
      result += `blockers: ${JSON.stringify(parsed.frontmatter.blockers)}\n`
    }
    if (parsed.frontmatter.blocks && parsed.frontmatter.blocks.length > 0) {
      result += `blocks: ${JSON.stringify(parsed.frontmatter.blocks)}\n`
    }
    if (parsed.frontmatter.parent) result += `parent: ${parsed.frontmatter.parent}\n`

    result += "---\n"
    if (parsed.body) result += parsed.body

    setContent(result)
  }

  const saveFile = async () => {
    if (!activePath()) return
    try {
      setStatus("Saving...")
      await fetchJSON(apiURL(`/api/file?path=${encodeURIComponent(activePath())}`), {
        method: "PUT",
        body: JSON.stringify({ content: content() }),
      })
      setDirty(false)
      setStatus(`Saved ${activePath()}`)
    } catch (err) {
      setStatus(`Save failed: ${errorMessage(err)}`)
    }
  }

  const onSelect = (entry: TaskTreeNode) => {
    const path = entry.task.path
    if (!path) return
    void loadFile(path)
  }

  const hasChildren = (node: TaskTreeNode) => (taskChildren.get(node.task.short_id)?.length ?? 0) > 0
  const isExpanded = (node: TaskTreeNode) => node.expanded ?? !node.task.completed

  const toggleNode = (node: TaskTreeNode) => {
    node.expanded = !isExpanded(node)
  }

  const availableRoles = createMemo(() => {
    const roles = new Set<string>()
    tasks().forEach((task) => {
      if (task.role) roles.add(task.role)
    })
    return Array.from(roles).sort()
  })

  const availablePriorities = createMemo(() => {
    const priorities = new Set<string>()
    tasks().forEach((task) => {
      if (task.priority) priorities.add(task.priority)
    })
    return Array.from(priorities).sort()
  })

  const taskNodesById = new ReactiveMap<string, TaskTreeNode>()
  const taskChildren = new ReactiveMap<string, string[]>()
  const sortedNodes = createMutable<TaskTreeNode[]>([])

  const taskCompare = (field: SortField, direction: "asc" | "desc") => (a: TaskItem, b: TaskItem) => {
    // if (field === "actionable") {
    //   const aBlocked = a.blockers && a.blockers.length > 0
    //   const bBlocked = b.blockers && b.blockers.length > 0
    //   if (aBlocked !== bBlocked) return aBlocked ? 1 : -1
    //   if (a.priority !== b.priority) {
    //     const priorityOrder: Record<string, number> = { high: 3, medium: 2, low: 1, "": 0 }
    //     const aPri = priorityOrder[a.priority] || 0
    //     const bPri = priorityOrder[b.priority] || 0
    //     return bPri - aPri
    //   }
    //   return b.date_edited.localeCompare(a.date_edited)
    // }

    let aVal: any = a[field]
    let bVal: any = b[field]

    if (field === "priority") {
      const priorityOrder: Record<string, number> = { high: 3, medium: 2, low: 1, "": 0 }
      aVal = priorityOrder[aVal] || 0
      bVal = priorityOrder[bVal] || 0
    }

    if (aVal < bVal) return direction === "asc" ? -1 : 1
    if (aVal > bVal) return direction === "asc" ? 1 : -1
    return 0
  }

  // re-sort sortedNodes when sortField/Direction changes
  createMemo(() => {
    const field = sortField()
    const direction = sortDirection()

    const taskCmp = taskCompare(field, direction)
    untrack(() => sortedNodes.sort((a, b) => taskCmp(a.task, b.task)))
  })


  let lastUpdateTick = 0
  createEffect(() => {
    const field = sortField()
    const taskCmp = untrack(() => taskCompare(field, sortDirection()))
    const nodeCmp = (a: TaskTreeNode, b: TaskTreeNode) => taskCmp(a.task, b.task)
    const stats = { added: 0, update: 0, removed: 0 }
    lastUpdateTick += 1
    for (const task of tasks()) {
      let item = untrack(() => taskNodesById.get(task.short_id))
      if (!item) {
        // ITEM ADD
        stats.added += 1
        item = createMutable({ task, depth: 0, lastUpdateTick })
        taskNodesById.set(task.short_id, item)
        if (task.parent) {
          let siblings = untrack(() => taskChildren.get(task!.parent))
          if (!siblings) {
            siblings = [task.short_id]
          } else {
            sortedInsert(siblings, task.short_id)
          }
          taskChildren.set(task!.parent, siblings)
        }
        sortedNodes.splice(sortedIndexCmp(sortedNodes, nodeCmp, item), 0, item)

      } else {
        // ITEM UPDATE
        stats.update += 1
        const removeThenAdd = item.task[field] != task[field]
        if (removeThenAdd) {
          // remove
          sortedNodes.splice(sortedIndexCmp(sortedNodes, nodeCmp, item), 1)
        }
        Object.assign(item.task, task)
        item.lastUpdateTick = lastUpdateTick
        if (removeThenAdd) {
          // add
          sortedNodes.splice(sortedIndexCmp(sortedNodes, nodeCmp, item), 0, item)
        }
      }
    }
    // now delete any old ones
    untrack(() => {
      for (const node of taskNodesById.values()) {
        if (node.lastUpdateTick != lastUpdateTick) {
          // ITEM REMOVE

          stats.removed += 1
          taskNodesById.delete(node.task.short_id)
          sortedNodes.splice(sortedIndexCmp(sortedNodes, nodeCmp, node), 1)
          let siblings = untrack(() => taskChildren.get(node.task.parent))
          if (siblings)
            sortedRemove(siblings, node.task.short_id)
        }
      }
    })

    // console.log('tasks update', stats)
  })

  const taskTreeFlattened = createMemo(() => {
    console.log('flattening...')
    const tasksSeen = new Set<string>()
    const flattened: TaskTreeNode[] = []

    const pushNodeFlattened = (node: TaskTreeNode, depth: number, expanded: boolean) => {
      tasksSeen.add(node.task.short_id)
      node.depth = depth
      if (expanded)
        flattened.push(node)
      expanded &&= isExpanded(node)
      const children = taskChildren.get(node.task.short_id)
      if (children) {
        for (const childId of children) {
          const child = taskNodesById.get(childId)
          if (child) {
            pushNodeFlattened(child, depth + 1, expanded)
          }
        }
      }
    }

    for (let node of sortedNodes) {
      if (tasksSeen.has(node.task.short_id))
        continue

      while (node.task.parent) {
        const parent = taskNodesById.get(node.task.parent)
        if (!parent) break
        node = parent
      }

      // now that the root task for this first task has been found, we expand it
      pushNodeFlattened(node, 0, true)
    }
    return flattened
  })

  const filteredTasks = createMemo(() => {
    let filtered = tasks()
    const query = searchQuery().toLowerCase()

    if (query) {
      filtered = filtered.filter((task) =>
        task.title.toLowerCase().includes(query) ||
        task.short_id.toLowerCase().includes(query) ||
        task.id.toLowerCase().includes(query) ||
        task.role.toLowerCase().includes(query)
      )
    }

    const status = filterStatus()
    if (status !== "all") {
      filtered = filtered.filter((task) => status === "done" ? task.completed : !task.completed)
    }

    const role = filterRole()
    if (role !== "all") {
      filtered = filtered.filter((task) => task.role === role)
    }

    const priority = filterPriority()
    if (priority !== "all") {
      filtered = filtered.filter((task) => task.priority === priority)
    }

    return filtered
  })

  const handleSortChange = (field: SortField) => {
    if (sortField() === field) {
      setSortDirection((d) => (d === "asc" ? "desc" : "asc"))
    } else {
      setSortField(field)
      setSortDirection(field === "title" || field === "short_id" ? "asc" : "desc")
    }
  }

  const handleParsedTaskChange = (task: ParsedTask) => {
    setParsedTask(task)
    updateContent()
    setDirty(true)
  }

  const handleContentChange = (newContent: string) => {
    setContent(newContent)
    setDirty(true)
  }

  createEffect(() => {
    const current = tab()
    setActivePath("")
    setContent("")
    setDirty(false)
    setStatus("")
    if (current === "tasks") {
      void loadTasks()
    }
  })

  createEffect(() => {
    const project = currentProject()
    if (!project) return
    setActivePath("")
    setContent("")
    setDirty(false)
    if (tab() === "tasks") {
      void loadTasks()
    }
  })

  onMount(() => {
    void loadProjects()
    void loadTasks()

    const source = new EventSource("/api/stream")
    const onOpen = () => setConnected(true)
    const onError = () => setConnected(false)
    const onTask = (event: MessageEvent) => {
      try {
        const update = JSON.parse(event.data) as StreamUpdate
        if (update.project !== currentProject()) return
        setLastEvent(`${update.event} • ${update.path}`)
        if (tab() === "tasks") void loadTasks()
        if (activePath() && update.path === activePath()) void loadFile(activePath())
      } catch (err) {
        setStatus(`Stream error: ${errorMessage(err)}`)
      }
    }

    source.addEventListener("open", onOpen)
    source.addEventListener("error", onError)
    source.addEventListener("task", onTask as EventListener)

    const keyHandler = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "s") {
        event.preventDefault()
        void saveFile()
      }
    }
    window.addEventListener("keydown", keyHandler)

    onCleanup(() => {
      source.removeEventListener("open", onOpen)
      source.removeEventListener("error", onError)
      source.removeEventListener("task", onTask as EventListener)
      source.close()
      window.removeEventListener("keydown", keyHandler)
    })
  })

  // const [testVisible, setTestVisible] = createSignal(true)

  return (
    <div class="app">
      {/* <button onClick={() => setTestVisible(x => !x)}>toggle</button>
      <table>
        <MyTransitionGroup classPrefix="tree-item">
          <tr id="one" class="tree-item"><td style="border: 1px solid gray;">One</td></tr>
          <Show when={testVisible()}>
            <tr id="two" class="tree-item"><td style="border: 1px solid gray;">Two</td></tr>
          </Show>
          <tr id="three" class="tree-item"><td style="border: 1px solid gray;">Three</td></tr>
        </MyTransitionGroup>
      </table> */}

      {/* <div>
        <MyTransitionGroup classPrefix="tree-item">
          <div id="one" class="tree-item" style="border: 1px solid gray;"><span>One</span></div>
          <Show when={testVisible()}>
            <div id="two" class="tree-item" style="border: 1px solid gray;"><span>Two</span></div>
          </Show>
          <div id="three" class="tree-item" style="border: 1px solid gray;"><span>Three</span></div>
        </MyTransitionGroup>
      </div> */}


      <Header
        projects={projects()}
        currentProject={currentProject()}
        connected={connected()}
        onProjectChange={setCurrentProject}
      />

      <section class="workspace">
        <Sidebar
          tab={tab()}
          projects={projects()}
          currentProject={currentProject()}
          searchQuery={searchQuery()}
          filterStatus={filterStatus()}
          filterRole={filterRole()}
          filterPriority={filterPriority()}
          availableRoles={availableRoles()}
          availablePriorities={availablePriorities()}
          onTabChange={setTab}
          onProjectChange={setCurrentProject}
          onSearchChange={setSearchQuery}
          onFilterStatusChange={setFilterStatus}
          onFilterRoleChange={setFilterRole}
          onFilterPriorityChange={setFilterPriority}
        />

        <div class="list-pane">
          <div class="pane-header">
            <h2>Tasks Library</h2>
            <span class="pill">{filteredTasks().length} items</span>
          </div>
          <div class="list">
            <TaskTable
              tasks={taskTreeFlattened()}
              activePath={activePath()}
              sortField={sortField()}
              sortDirection={sortDirection()}
              hasChildren={hasChildren}
              isExpanded={isExpanded}
              onSelect={onSelect}
              onToggleNode={toggleNode}
              onSortChange={handleSortChange}
            />
          </div>
        </div>

        <Editor
          activePath={activePath()}
          content={content()}
          dirty={dirty()}
          status={status()}
          lastEvent={lastEvent()}
          tab={tab()}
          parsedTask={parsedTask()}
          onContentChange={handleContentChange}
          onParsedTaskChange={handleParsedTaskChange}
          onSave={saveFile}
        />
      </section>
    </div>
  )
}
