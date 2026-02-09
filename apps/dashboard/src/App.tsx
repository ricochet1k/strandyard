import { createEffect, createResource, createSignal, onCleanup, onMount, createMemo, untrack, Show } from "solid-js"
import Header from "./components/Header"
import Sidebar from "./components/Sidebar"
import TaskTable from "./components/TaskTable"
import RolesTable from "./components/RolesTable"
import TemplatesTable from "./components/TemplatesTable"
import Editor from "./components/Editor"
import AddTaskModal from "./components/AddTaskModal"
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
  status: string
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

export type RoleItem = {
  name: string
  path: string
  description: string
}

export type RoleDetail = RoleItem & {
  body: string
}

export type TemplateItem = {
  name: string
  path: string
  role: string
  priority: string
  description: string
  id_prefix: string
}

export type TemplateDetail = TemplateItem & {
  body: string
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

type TaskDetail = {
  id: string
  short_id: string
  title: string
  role: string
  priority: string
  completed: boolean
  status: string
  parent: string
  blockers: string[]
  blocks: string[]
  path: string
  date_created: string
  date_edited: string
  body: string
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

function buildProjectUrl(path: string, project: string, configure?: (params: URLSearchParams) => void) {
  const params = new URLSearchParams()
  params.set("project", project)
  configure?.(params)
  return `${path}?${params.toString()}`
}

async function fetchTasksForProject(project: string | undefined) {
  if (!project) {
    return []
  }
  return fetchJSON<TaskItem[]>(buildProjectUrl("/api/tasks", project))
}

async function fetchRolesForProject(project: string | undefined) {
  if (!project) {
    return []
  }
  return await fetchJSON<RoleItem[]>(buildProjectUrl("/api/roles", project))
}

async function fetchTemplatesForProject(project: string | undefined) {
  if (!project) {
    return []
  }
  return await fetchJSON<TemplateItem[]>(buildProjectUrl("/api/templates", project))
}

function errorMessage(err: unknown) {
  if (err instanceof Error) {
    return err.message
  }
  return "Unknown error"
}

function normalizeTaskStatus(task: { status: string; completed: boolean }) {
  const raw = task.status?.trim()
  if (raw) return raw
  return task.completed ? "done" : "open"
}

function isActiveStatus(status: string) {
  return status === "open" || status === "in_progress"
}

export default function App() {
  const [tab, setTab] = createSignal<Tab>("tasks")
  const [activeTaskDetail, setActiveTaskDetail] = createSignal<TaskDetail | null>(null)
  const [activeRoleDetail, setActiveRoleDetail] = createSignal<RoleDetail | null>(null)
  const [activeTemplateDetail, setActiveTemplateDetail] = createSignal<TemplateDetail | null>(null)
  const [dirty, setDirty] = createSignal(false)
  const [status, setStatus] = createSignal("")
  const [connected, setConnected] = createSignal(false)
  const [lastEvent, setLastEvent] = createSignal("")
  const [projects, setProjects] = createSignal<ProjectInfo[]>([])
  const [currentProject, setCurrentProject] = createSignal("")
  const [showAddTaskModal, setShowAddTaskModal] = createSignal(false)
  const [addTaskParent, setAddTaskParent] = createSignal("")
  const [originId, setOriginId] = createSignal<string | null>(null)
  const [relationship, setRelationship] = createSignal<string | null>(null)

  const [tasks, { refetch: reloadTasks }] = createResource(currentProject, fetchTasksForProject, { initialValue: [] })
  const [roles, { refetch: reloadRoles }] = createResource(currentProject, fetchRolesForProject, { initialValue: [] })
  const [templates, { refetch: reloadTemplates }] = createResource(currentProject, fetchTemplatesForProject, { initialValue: [] })

  createEffect(() => {
    const err = tasks.error
    if (err) {
      setStatus(`Failed to load tasks: ${errorMessage(err)}`)
    }
  })

  createEffect(() => {
    const err = roles.error
    if (err) {
      setStatus(`Failed to load roles: ${errorMessage(err)}`)
    }
  })

  createEffect(() => {
    const err = templates.error
    if (err) {
      setStatus(`Failed to load templates: ${errorMessage(err)}`)
    }
  })

  // Task filtering, sorting, and search
  const [searchQuery, setSearchQuery] = createSignal("")
  const [filterStatus, setFilterStatus] = createSignal<
    "all" | "active" | "open" | "in_progress" | "done" | "cancelled" | "duplicate"
  >("active")
  const [filterRole, setFilterRole] = createSignal<string>("all")
  const [filterPriority, setFilterPriority] = createSignal<string>("all")
  const [hideBlocked, setHideBlocked] = createSignal(false)
  const [viewMode, setViewMode] = createSignal<"tree" | "list">("tree")
  const [sortField, setSortField] = createSignal<SortField>("priority")
  const [sortDirection, setSortDirection] = createSignal<SortDirection>("desc")

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

  const loadRole = async (path: string) => {
    try {
      const data = await fetchJSON<RoleDetail>(apiURL(`/api/role?path=${encodeURIComponent(path)}`))
      setActiveRoleDetail(data)
      setDirty(false)
      setStatus(`Loaded ${path}`)
    } catch (err) {
      setStatus(`Failed to load role: ${errorMessage(err)}`)
    }
  }

  const loadTemplate = async (path: string) => {
    try {
      const data = await fetchJSON<TemplateDetail>(apiURL(`/api/template?path=${encodeURIComponent(path)}`))
      setActiveTemplateDetail(data)
      setDirty(false)
      setStatus(`Loaded ${path}`)
    } catch (err) {
      setStatus(`Failed to load template: ${errorMessage(err)}`)
    }
  }

  const saveRole = async () => {
    const role = activeRoleDetail()
    if (!role) return
    try {
      setStatus("Saving...")
      const updated = await fetchJSON<RoleDetail>(apiURL(`/api/role?path=${encodeURIComponent(role.path)}`), {
        method: "PUT",
        body: JSON.stringify({
          description: role.description,
          body: role.body,
        }),
      })
      setActiveRoleDetail(updated)
      setDirty(false)
      setStatus(`Saved ${role.path}`)
      // Reload roles to update the list
      await reloadRoles()
    } catch (err) {
      setStatus(`Save failed: ${errorMessage(err)}`)
    }
  }

  const saveTemplate = async () => {
    const template = activeTemplateDetail()
    if (!template) return
    try {
      setStatus("Saving...")
      const updated = await fetchJSON<TemplateDetail>(apiURL(`/api/template?path=${encodeURIComponent(template.path)}`), {
        method: "PUT",
        body: JSON.stringify({
          role: template.role,
          priority: template.priority,
          description: template.description,
          id_prefix: template.id_prefix,
          body: template.body,
        }),
      })
      setActiveTemplateDetail(updated)
      setDirty(false)
      setStatus(`Saved ${template.path}`)
      // Reload templates to update the list
      await reloadTemplates()
    } catch (err) {
      setStatus(`Save failed: ${errorMessage(err)}`)
    }
  }

  const loadTask = async (taskId: string, rel?: string, orig?: string) => {
    try {
      const data = await fetchJSON<TaskDetail>(apiURL(`/api/task?id=${encodeURIComponent(taskId)}`))
      setActiveTaskDetail(data)
      setDirty(false)
      setStatus(`Loaded ${data.short_id}`)
      setRelationship(rel || null)
      setOriginId(orig || null)

      // Update URL
      const params = new URLSearchParams(window.location.search)
      params.set("task", taskId)
      if (rel) params.set("relationship", rel)
      else params.delete("relationship")
      if (orig) params.set("origin", orig)
      else params.delete("origin")

      const newUrl = `${window.location.pathname}?${params.toString()}`
      window.history.pushState({ taskId, rel, orig }, "", newUrl)
    } catch (err) {
      setStatus(`Failed to load task: ${errorMessage(err)}`)
    }
  }

  const onSelectTask = (id: string, rel?: string, orig?: string) => {
    void loadTask(id, rel, orig)
  }

  const saveTask = async () => {
    const task = activeTaskDetail()
    if (!task) return
    try {
      setStatus("Saving...")
      const updated = await fetchJSON<TaskDetail>(apiURL(`/api/task?id=${encodeURIComponent(task.id)}`), {
        method: "PATCH",
        body: JSON.stringify({
          title: task.title,
          role: task.role,
          priority: task.priority,
          status: task.status,
          parent: task.parent,
          blockers: task.blockers,
          blocks: task.blocks,
          body: task.body,
        }),
      })
      setActiveTaskDetail(updated)
      setDirty(false)
      setStatus(`Saved ${task.short_id}`)
    } catch (err) {
      setStatus(`Save failed: ${errorMessage(err)}`)
    }
  }

  const onSelect = (entry: TaskTreeNode) => {
    void loadTask(entry.task.id)
  }

  const handlePopState = (event: PopStateEvent) => {
    const state = event.state
    if (state && state.taskId) {
      void loadTask(state.taskId, state.rel, state.orig)
    }
  }

  const onSelectRole = (role: RoleItem) => {
    void loadRole(role.path)
  }

  const onSelectTemplate = (template: TemplateItem) => {
    void loadTemplate(template.path)
  }

  const handleRoleDetailChange = (updated: RoleDetail) => {
    setActiveRoleDetail(updated)
    setDirty(true)
  }

  const handleTemplateDetailChange = (updated: TemplateDetail) => {
    setActiveTemplateDetail(updated)
    setDirty(true)
  }

  const handleAddTask = async (data: {
    template_name: string
    title: string
    role: string
    priority: string
    parent: string
    body: string
  }) => {
    try {
      setStatus("Creating task...")
      await fetchJSON(apiURL("/api/task"), {
        method: "POST",
        body: JSON.stringify(data),
      })
      setStatus("Task created successfully")
      // Reload tasks to show the new task
      await reloadTasks()
    } catch (err) {
      setStatus(`Failed to create task: ${errorMessage(err)}`)
      throw err
    }
  }

  const openAddTaskModal = (parent = "") => {
    console.log("Opening add task modal with parent:", parent)
    setAddTaskParent(parent)
    setShowAddTaskModal(true)
  }

  const handleAddSubtask = () => {
    const currentTask = activeTaskDetail()
    if (!currentTask) return

    // Open the add task modal with the current task as parent
    openAddTaskModal(currentTask.id)
  }

  const hasChildren = (node: TaskTreeNode) => (taskChildren.get(node.task.short_id)?.length ?? 0) > 0
  const isExpanded = (node: TaskTreeNode) => node.expanded ?? isActiveStatus(normalizeTaskStatus(node.task))

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

  // Map to convert short_id to short_id for parent lookups
  const shortIdMap = new ReactiveMap<string, string>()

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

    // First pass: build short_id map
    for (const task of tasks()) {
      shortIdMap.set(task.short_id, task.short_id)
      shortIdMap.set(task.id, task.short_id)
    }

    for (const task of tasks()) {
      let item = untrack(() => taskNodesById.get(task.short_id))
      if (!item) {
        // ITEM ADD
        stats.added += 1
        item = createMutable({ task, depth: 0, lastUpdateTick })
        taskNodesById.set(task.short_id, item)
        if (task.parent) {
          // Parent field could be either short_id or full ID, normalize it
          const parentShortId = shortIdMap.get(task.parent) || task.parent
          let siblings = untrack(() => taskChildren.get(parentShortId))
          if (!siblings) {
            siblings = [task.short_id]
          } else {
            sortedInsert(siblings, task.short_id)
          }
          taskChildren.set(parentShortId, siblings)
        }
        sortedNodes.splice(sortedIndexCmp(sortedNodes, nodeCmp, item), 0, item)

      } else {
        // ITEM UPDATE
        stats.update += 1
        const oldParent = item.task.parent
        const removeThenAdd = item.task[field] != task[field]
        if (removeThenAdd) {
          // remove
          sortedNodes.splice(sortedIndexCmp(sortedNodes, nodeCmp, item), 1)
        }
        Object.assign(item.task, task)
        item.lastUpdateTick = lastUpdateTick

        // Handle parent change
        if (oldParent !== task.parent) {
          // Remove from old parent
          if (oldParent) {
            const oldParentShortId = shortIdMap.get(oldParent) || oldParent
            let oldSiblings = untrack(() => taskChildren.get(oldParentShortId))
            if (oldSiblings) {
              sortedRemove(oldSiblings, task.short_id)
            }
          }
          // Add to new parent
          if (task.parent) {
            const newParentShortId = shortIdMap.get(task.parent) || task.parent
            let siblings = untrack(() => taskChildren.get(newParentShortId))
            if (!siblings) {
              siblings = [task.short_id]
              taskChildren.set(newParentShortId, siblings)
            } else {
              sortedInsert(siblings, task.short_id)
            }
          }
        }

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
          if (node.task.parent) {
            const parentShortId = shortIdMap.get(node.task.parent) || node.task.parent
            let siblings = untrack(() => taskChildren.get(parentShortId))
            if (siblings)
              sortedRemove(siblings, node.task.short_id)
          }
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

  const matchesFilters = (task: TaskItem) => {
    const query = searchQuery().toLowerCase()
    if (query) {
      const matches = task.title.toLowerCase().includes(query) ||
        task.short_id.toLowerCase().includes(query) ||
        task.id.toLowerCase().includes(query) ||
        task.role.toLowerCase().includes(query)
      if (!matches) return false
    }

    const status = filterStatus()
    if (status !== "all") {
      const taskStatus = normalizeTaskStatus(task)
      if (status === "active") {
        if (!isActiveStatus(taskStatus)) return false
      } else if (taskStatus !== status) {
        return false
      }
    }

    const role = filterRole()
    if (role !== "all") {
      if (task.role !== role) return false
    }

    const priority = filterPriority()
    if (priority !== "all") {
      if (task.priority !== priority) return false
    }

    if (hideBlocked()) {
      if (task.blockers && task.blockers.length > 0) return false
    }

    return true
  }

  const filteredTaskNodes = createMemo(() => {
    const mode = viewMode()
    const flattened = taskTreeFlattened()

    if (mode === "list") {
      // List view: filter individual tasks
      return flattened.filter(node => matchesFilters(node.task))
    } else {
      // Tree view: keep hierarchy, but filter out non-matching branches
      const matchingIds = new Set<string>()
      const ancestorIds = new Set<string>()

      // First pass: find all matching tasks
      for (const node of flattened) {
        if (matchesFilters(node.task)) {
          matchingIds.add(node.task.short_id)

          // Mark all ancestors as needed
          let parentId = node.task.parent
          while (parentId) {
            const parentShortId = shortIdMap.get(parentId) || parentId
            ancestorIds.add(parentShortId)
            const parent = taskNodesById.get(parentShortId)
            if (!parent) break
            parentId = parent.task.parent
          }
        }
      }

      // Second pass: include matching tasks and their ancestors
      return flattened.filter(node =>
        matchingIds.has(node.task.short_id) || ancestorIds.has(node.task.short_id)
      )
    }
  })

  const handleSortChange = (field: SortField) => {
    if (sortField() === field) {
      setSortDirection((d) => (d === "asc" ? "desc" : "asc"))
    } else {
      setSortField(field)
      setSortDirection(field === "title" || field === "short_id" ? "asc" : "desc")
    }
  }

  const handleTaskDetailChange = (updated: TaskDetail) => {
    setActiveTaskDetail(updated)
    setDirty(true)
  }

  createEffect(() => {
    const current = tab()
    setActiveTaskDetail(null)
    setActiveRoleDetail(null)
    setActiveTemplateDetail(null)
    setOriginId(null)
    setRelationship(null)
    setDirty(false)
    setStatus("")
  })

  onMount(() => {
    void loadProjects().then(() => {
      const params = new URLSearchParams(window.location.search)
      const taskId = params.get('task')
      if (taskId) {
        void loadTask(taskId, params.get('relationship') || undefined, params.get('origin') || undefined)
      }
    })

    window.addEventListener("popstate", handlePopState)

    let source: EventSource | null = null
    let socket: WebSocket | null = null

    const onUpdate = (data: any) => {
      try {
        const update = typeof data === 'string' ? JSON.parse(data) : data
        if (update.project !== currentProject()) return
        setLastEvent(`${update.event} • ${update.path}`)
        if (tab() === "tasks") void reloadTasks()
        if (tab() === "roles" || update.path.includes(".strand/roles/")) void reloadRoles()
        if (tab() === "templates" || update.path.includes("templates/")) void reloadTemplates()
        const active = activeTaskDetail()
        if (active) {
          const updatedId = update.task?.id
          if (updatedId && updatedId === active.id) {
            void loadTask(active.id)
          } else if (update.path === active.path) {
            void loadTask(active.id)
          }
        }
      } catch (err) {
        setStatus(`Stream error: ${errorMessage(err)}`)
      }
    }

    const connectWS = () => {
      const wsUrl = new URL("/api/ws", window.location.href)
      wsUrl.protocol = wsUrl.protocol === "https:" ? "wss:" : "ws:"
      const params = new URLSearchParams(window.location.search)
      const token = params.get("token")
      if (token) {
        wsUrl.searchParams.set("token", token)
      }

      socket = new WebSocket(wsUrl.toString())
      socket.onopen = () => {
        setConnected(true)
        console.log("Websocket connected")
      }
      socket.onclose = () => {
        setConnected(false)
        console.log("Websocket closed, falling back to SSE")
        connectSSE()
      }
      socket.onmessage = (event) => {
        const data = JSON.parse(event.data)
        if (data.ping || data.status === "connected") return
        onUpdate(data)
      }
      socket.onerror = (err) => {
        console.error("Websocket error:", err)
        socket?.close()
      }
    }

    const connectSSE = () => {
      if (connected()) return
      source = new EventSource("/api/stream")
      source.onopen = () => setConnected(true)
      source.onerror = () => {
        setConnected(false)
        source?.close()
      }
      source.addEventListener("task", (event: MessageEvent) => {
        onUpdate(event.data)
      })
    }

    // Try Websocket first
    connectWS()

    const keyHandler = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "s") {
        event.preventDefault()
        void saveTask()
      }
    }
    window.addEventListener("keydown", keyHandler)

    onCleanup(() => {
      socket?.close()
      source?.close()
      window.removeEventListener("keydown", keyHandler)
      window.removeEventListener("popstate", handlePopState)
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
          hideBlocked={hideBlocked()}
          viewMode={viewMode()}
          availableRoles={availableRoles()}
          availablePriorities={availablePriorities()}
          onTabChange={setTab}
          onProjectChange={setCurrentProject}
          onSearchChange={setSearchQuery}
          onFilterStatusChange={setFilterStatus}
          onFilterRoleChange={setFilterRole}
          onFilterPriorityChange={setFilterPriority}
          onHideBlockedChange={setHideBlocked}
          onViewModeChange={setViewMode}
        />

        <div class="list-pane">
          <Show when={tab() === "tasks"}>
            <div class="pane-header">
              <h2>Tasks Library</h2>
              <div style={{ display: "flex", gap: "0.5rem", "align-items": "center" }}>
                <span class="pill">{filteredTaskNodes().length} items</span>
                <button class="button button-primary" onClick={() => openAddTaskModal("")}>
                  + Add Task
                </button>
              </div>
            </div>
            <div class="list">
              <TaskTable
                tasks={filteredTaskNodes()}
                activePath={activeTaskDetail()?.path ?? ""}
                sortField={sortField()}
                sortDirection={sortDirection()}
                viewMode={viewMode()}
                hasChildren={hasChildren}
                isExpanded={isExpanded}
                onSelect={onSelect}
                onToggleNode={toggleNode}
                onSortChange={handleSortChange}
              />
            </div>
          </Show>

          <Show when={tab() === "roles"}>
            <div class="pane-header">
              <h2>Roles</h2>
              <span class="pill">{roles().length} items</span>
            </div>
            <div class="list">
              <RolesTable
                roles={roles()}
                activePath={activeRoleDetail()?.path ?? ""}
                onSelect={onSelectRole}
              />
            </div>
          </Show>

          <Show when={tab() === "templates"}>
            <div class="pane-header">
              <h2>Templates</h2>
              <span class="pill">{templates().length} items</span>
            </div>
            <div class="list">
              <TemplatesTable
                templates={templates()}
                activePath={activeTemplateDetail()?.path ?? ""}
                onSelect={onSelectTemplate}
              />
            </div>
          </Show>
        </div>

        <Editor
          task={activeTaskDetail()}
          tasks={tasks()}
          role={activeRoleDetail()}
          template={activeTemplateDetail()}
          dirty={dirty()}
          status={status()}
          lastEvent={lastEvent()}
          tab={tab()}
          originId={originId()}
          relationship={relationship()}
          onTaskChange={handleTaskDetailChange}
          onRoleChange={handleRoleDetailChange}
          onTemplateChange={handleTemplateDetailChange}
          onSave={tab() === "tasks" ? saveTask : tab() === "roles" ? saveRole : saveTemplate}
          onAddSubtask={handleAddSubtask}
          onSelectTask={onSelectTask}
        />
      </section>

      <Show when={showAddTaskModal()}>
        <AddTaskModal
          templates={templates()}
          roles={roles()}
          tasks={tasks()}
          defaultParent={addTaskParent()}
          onClose={() => setShowAddTaskModal(false)}
          onSubmit={handleAddTask}
        />
      </Show>
    </div>
  )
}
