import { Show, For, createSignal, onCleanup, createEffect } from "solid-js"
import type { RoleDetail, TemplateDetail } from "../App"
import { EditorView, basicSetup } from "codemirror"
import { markdown } from "@codemirror/lang-markdown"
import { oneDark } from "@codemirror/theme-one-dark"
import { EditorState } from "@codemirror/state"
import "./Editor.css"

type TaskDetail = {
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
  body: string
}

type EditorProps = {
  task: TaskDetail | null
  tasks: Array<{ id: string; short_id: string; title: string }>
  role: RoleDetail | null
  template: TemplateDetail | null
  dirty: boolean
  status: string
  lastEvent: string
  tab: string
  originId?: string | null
  relationship?: string | null
  onTaskChange: (task: TaskDetail) => void
  onRoleChange: (role: RoleDetail) => void
  onTemplateChange: (template: TemplateDetail) => void
  onSave: () => void
  onAddSubtask?: () => void
  onSelectTask?: (id: string, rel?: string, orig?: string) => void
}

export default function Editor(props: EditorProps) {
  const [newBlocker, setNewBlocker] = createSignal("")
  const [newBlock, setNewBlock] = createSignal("")

  const getTaskTitle = (id: string) => {
    if (!props.tasks) return id
    const t = props.tasks.find(t => t.id === id || t.short_id === id)
    return t ? `${t.short_id} - ${t.title}` : id
  }

  let taskEditorContainer: HTMLDivElement | undefined
  let roleEditorContainer: HTMLDivElement | undefined
  let templateEditorContainer: HTMLDivElement | undefined

  let taskEditorView: EditorView | undefined
  let roleEditorView: EditorView | undefined
  let templateEditorView: EditorView | undefined

  const updateTaskField = (field: keyof TaskDetail, value: any) => {
    if (!props.task) return
    props.onTaskChange({ ...props.task, [field]: value })
  }

  const addBlocker = () => {
    const value = newBlocker().trim()
    if (!value) return
    const current = props.task?.blockers || []
    if (!current.includes(value)) {
      updateTaskField("blockers", [...current, value])
    }
    setNewBlocker("")
  }

  const removeBlocker = (id: string) => {
    const current = props.task?.blockers || []
    updateTaskField("blockers", current.filter((b) => b !== id))
  }

  const addBlock = () => {
    const value = newBlock().trim()
    if (!value) return
    const current = props.task?.blocks || []
    if (!current.includes(value)) {
      updateTaskField("blocks", [...current, value])
    }
    setNewBlock("")
  }

  const removeBlock = (id: string) => {
    const current = props.task?.blocks || []
    updateTaskField("blocks", current.filter((b) => b !== id))
  }

  const task = () => props.task
  const role = () => props.role
  const template = () => props.template

  const displayPath = () => {
    if (props.tab === "tasks") {
      return task()?.path || "Select a file to begin."
    } else if (props.tab === "roles") {
      return role()?.path || "Select a file to begin."
    } else if (props.tab === "templates") {
      return template()?.path || "Select a file to begin."
    }
    return "Select a file to begin."
  }

  const canSave = () => {
    if (props.tab === "tasks") return !!task()
    if (props.tab === "roles") return !!role()
    if (props.tab === "templates") return !!template()
    return false
  }

  const updateRoleField = (field: keyof RoleDetail, value: any) => {
    if (!role()) return
    props.onRoleChange({ ...role()!, [field]: value })
  }

  const updateTemplateField = (field: keyof TemplateDetail, value: any) => {
    if (!template()) return
    props.onTemplateChange({ ...template()!, [field]: value })
  }

  // Setup CodeMirror for task editor
  createEffect(() => {
    if (props.tab === "tasks" && task() && taskEditorContainer) {
      if (!taskEditorView) {
        const state = EditorState.create({
          doc: task()?.body || "",
          extensions: [
            basicSetup,
            markdown(),
            oneDark,
            EditorView.lineWrapping,
            EditorView.updateListener.of((update) => {
              if (update.docChanged && task()) {
                const newContent = update.state.doc.toString()
                props.onTaskChange({ ...task()!, body: newContent })
              }
            })
          ]
        })
        taskEditorView = new EditorView({
          state,
          parent: taskEditorContainer
        })
      } else {
        // Update editor content if task changed
        const currentContent = taskEditorView.state.doc.toString()
        const newContent = task()?.body || ""
        if (currentContent !== newContent) {
          taskEditorView.dispatch({
            changes: { from: 0, to: currentContent.length, insert: newContent }
          })
        }
      }
    } else if (taskEditorView && (!task() || props.tab !== "tasks")) {
      taskEditorView.destroy()
      taskEditorView = undefined
    }
  })

  // Setup CodeMirror for role editor
  createEffect(() => {
    if (props.tab === "roles" && role() && roleEditorContainer) {
      if (!roleEditorView) {
        const state = EditorState.create({
          doc: role()?.body || "",
          extensions: [
            basicSetup,
            markdown(),
            oneDark,
            EditorView.lineWrapping,
            EditorView.updateListener.of((update) => {
              if (update.docChanged && role()) {
                const newContent = update.state.doc.toString()
                props.onRoleChange({ ...role()!, body: newContent })
              }
            })
          ]
        })
        roleEditorView = new EditorView({
          state,
          parent: roleEditorContainer
        })
      } else {
        const currentContent = roleEditorView.state.doc.toString()
        const newContent = role()?.body || ""
        if (currentContent !== newContent) {
          roleEditorView.dispatch({
            changes: { from: 0, to: currentContent.length, insert: newContent }
          })
        }
      }
    } else if (roleEditorView && (!role() || props.tab !== "roles")) {
      roleEditorView.destroy()
      roleEditorView = undefined
    }
  })

  // Setup CodeMirror for template editor
  createEffect(() => {
    if (props.tab === "templates" && template() && templateEditorContainer) {
      if (!templateEditorView) {
        const state = EditorState.create({
          doc: template()?.body || "",
          extensions: [
            basicSetup,
            markdown(),
            oneDark,
            EditorView.lineWrapping,
            EditorView.updateListener.of((update) => {
              if (update.docChanged && template()) {
                const newContent = update.state.doc.toString()
                props.onTemplateChange({ ...template()!, body: newContent })
              }
            })
          ]
        })
        templateEditorView = new EditorView({
          state,
          parent: templateEditorContainer
        })
      } else {
        const currentContent = templateEditorView.state.doc.toString()
        const newContent = template()?.body || ""
        if (currentContent !== newContent) {
          templateEditorView.dispatch({
            changes: { from: 0, to: currentContent.length, insert: newContent }
          })
        }
      }
    } else if (templateEditorView && (!template() || props.tab !== "templates")) {
      templateEditorView.destroy()
      templateEditorView = undefined
    }
  })

  onCleanup(() => {
    if (taskEditorView) taskEditorView.destroy()
    if (roleEditorView) roleEditorView.destroy()
    if (templateEditorView) templateEditorView.destroy()
  })

  return (
    <div class="editor-pane">
      <div class="pane-header">
        <div>
          <h2>Editor</h2>
          <p class="detail">{displayPath()}</p>
        </div>
        <div class="editor-actions">
          <Show when={props.tab === "tasks" && task() && props.onAddSubtask}>
            <button class="primary" onClick={() => props.onAddSubtask?.()}>
              + Add Subtask
            </button>
          </Show>
          <span class={`sync ${props.dirty ? "dirty" : ""}`}>
            {props.dirty ? "Unsaved" : "Synced"}
          </span>
          <button class="primary" disabled={!canSave()} onClick={() => props.onSave()}>
            Save
          </button>
        </div>
      </div>

      <Show when={props.tab === "tasks" && task()}>
        <div class="editor-container">
          <Show when={props.originId}>
            <div class="origin-header">
              Arrived from:{" "}
              <button class="tag-link" onClick={() => props.onSelectTask?.(props.originId!)}>
                {getTaskTitle(props.originId!)}
              </button>
              <span class="relationship-label">({props.relationship})</span>
            </div>
          </Show>
          <div class="editor-controls">
            <div class="editor-field">
              <label for="task-title">Title</label>
              <input
                id="task-title"
                type="text"
                class="editor-input"
                value={task()?.title || ""}
                onInput={(e) => updateTaskField("title", e.currentTarget.value)}
              />
            </div>

            <div class="editor-field">
              <label for="task-parent">Parent</label>
              <select
                id="task-parent"
                class="editor-input"
                value={task()?.parent || ""}
                onChange={(e) => updateTaskField("parent", e.currentTarget.value)}
              >
                <option value="">No Parent</option>
                <For each={props.tasks}>
                  {(t) => (
                    <Show when={t.id !== task()?.id}>
                      <option value={t.id}>
                        {t.short_id} - {t.title}
                      </option>
                    </Show>
                  )}
                </For>
              </select>
            </div>

            <div class="editor-field">
              <label for="task-role">Role</label>
              <input
                id="task-role"
                type="text"
                class="editor-input"
                value={task()?.role || ""}
                onInput={(e) => updateTaskField("role", e.currentTarget.value)}
              />
            </div>

            <div class="editor-field">
              <label for="task-priority">Priority</label>
              <select
                id="task-priority"
                class="editor-input"
                value={task()?.priority || ""}
                onChange={(e) => updateTaskField("priority", e.currentTarget.value)}
              >
                <option value="">None</option>
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
              </select>
            </div>

            <div class="editor-field">
              <label for="task-status">Status</label>
              <select
                id="task-status"
                class="editor-input"
                value={task()?.completed ? "done" : "active"}
                onChange={(e) => updateTaskField("completed", e.currentTarget.value === "done")}
              >
                <option value="active">Active</option>
                <option value="done">Done</option>
              </select>
            </div>

            <div class="editor-field">
              <label>Blockers (blocked by)</label>
              <div class="tag-list">
                <For each={task()?.blockers || []}>
                  {(blocker) => (
                    <span class={`tag ${props.originId === blocker && props.relationship === "blocked-by" ? "origin" : ""}`}>
                      <button class="tag-link" onClick={() => props.onSelectTask?.(blocker, "blocking", task()?.id)}>
                        {getTaskTitle(blocker)}
                      </button>
                      <button class="tag-remove" onClick={() => removeBlocker(blocker)}>×</button>
                    </span>
                  )}
                </For>
              </div>
              <div class="tag-input-row">
                <input
                  type="text"
                  class="editor-input"
                  placeholder="Task ID..."
                  value={newBlocker()}
                  onInput={(e) => setNewBlocker(e.currentTarget.value)}
                  onKeyDown={(e) => e.key === "Enter" && addBlocker()}
                />
                <button class="add-btn" onClick={addBlocker}>Add</button>
              </div>
            </div>

            <div class="editor-field">
              <label>Blocks (blocking)</label>
              <div class="tag-list">
                <For each={task()?.blocks || []}>
                  {(block) => (
                    <span class={`tag ${props.originId === block && props.relationship === "blocking" ? "origin" : ""}`}>
                      <button class="tag-link" onClick={() => props.onSelectTask?.(block, "blocked-by", task()?.id)}>
                        {getTaskTitle(block)}
                      </button>
                      <button class="tag-remove" onClick={() => removeBlock(block)}>×</button>
                    </span>
                  )}
                </For>
              </div>
              <div class="tag-input-row">
                <input
                  type="text"
                  class="editor-input"
                  placeholder="Task ID..."
                  value={newBlock()}
                  onInput={(e) => setNewBlock(e.currentTarget.value)}
                  onKeyDown={(e) => e.key === "Enter" && addBlock()}
                />
                <button class="add-btn" onClick={addBlock}>Add</button>
              </div>
            </div>
          </div>

          <div class="editor-body codemirror-container" ref={taskEditorContainer} />
        </div>
      </Show>

      <Show when={props.tab === "roles" && role()}>
        <div class="editor-container">
          <div class="editor-controls">
            <div class="editor-field">
              <label for="role-description">Description</label>
              <input
                id="role-description"
                type="text"
                class="editor-input"
                value={role()?.description || ""}
                onInput={(e) => updateRoleField("description", e.currentTarget.value)}
              />
            </div>
          </div>

          <div class="editor-body codemirror-container" ref={roleEditorContainer} />
        </div>
      </Show>

      <Show when={props.tab === "templates" && template()}>
        <div class="editor-container">
          <div class="editor-controls">
            <div class="editor-field">
              <label for="template-role">Role</label>
              <input
                id="template-role"
                type="text"
                class="editor-input"
                value={template()?.role || ""}
                onInput={(e) => updateTemplateField("role", e.currentTarget.value)}
              />
            </div>

            <div class="editor-field">
              <label for="template-priority">Priority</label>
              <select
                id="template-priority"
                class="editor-input"
                value={template()?.priority || ""}
                onChange={(e) => updateTemplateField("priority", e.currentTarget.value)}
              >
                <option value="">None</option>
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
              </select>
            </div>

            <div class="editor-field">
              <label for="template-description">Description</label>
              <input
                id="template-description"
                type="text"
                class="editor-input"
                value={template()?.description || ""}
                onInput={(e) => updateTemplateField("description", e.currentTarget.value)}
              />
            </div>

            <div class="editor-field">
              <label for="template-id-prefix">ID Prefix</label>
              <input
                id="template-id-prefix"
                type="text"
                class="editor-input"
                value={template()?.id_prefix || ""}
                onInput={(e) => updateTemplateField("id_prefix", e.currentTarget.value)}
              />
            </div>
          </div>

          <div class="editor-body codemirror-container" ref={templateEditorContainer} />
        </div>
      </Show>

      <Show when={!task() && !role() && !template()}>
        <textarea
          class="editor"
          value=""
          placeholder="Select a task, role, or template to edit."
          disabled
        />
      </Show>

      <div class="footer">
        <span>{props.status}</span>
        <span>{props.lastEvent}</span>
      </div>
    </div>
  )
}
