import { Show, For, createSignal } from "solid-js"
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
  dirty: boolean
  status: string
  lastEvent: string
  tab: string
  onTaskChange: (task: TaskDetail) => void
  onSave: () => void
}

export default function Editor(props: EditorProps) {
  const [newBlocker, setNewBlocker] = createSignal("")
  const [newBlock, setNewBlock] = createSignal("")

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

  return (
    <div class="editor-pane">
      <div class="pane-header">
        <div>
          <h2>Editor</h2>
          <p class="detail">
            {task()?.path ? task()?.path : "Select a file to begin."}
          </p>
        </div>
        <div class="editor-actions">
          <span class={`sync ${props.dirty ? "dirty" : ""}`}>
            {props.dirty ? "Unsaved" : "Synced"}
          </span>
          <button class="primary" disabled={!task()} onClick={() => props.onSave()}>
            Save
          </button>
        </div>
      </div>

      <Show
        when={props.tab === "tasks" && task()}
        fallback={
          <textarea
            class="editor"
            value=""
            placeholder="Select a task, role, or template file to edit."
          />
        }
      >
        <div class="editor-container">
          <div class="editor-controls">
            <div class="editor-field">
              <label>Title</label>
              <input
                type="text"
                class="editor-input"
                value={task()?.title || ""}
                onInput={(e) => updateTaskField("title", e.currentTarget.value)}
              />
            </div>

            <div class="editor-field">
              <label>Role</label>
              <input
                type="text"
                class="editor-input"
                value={task()?.role || ""}
                onInput={(e) => updateTaskField("role", e.currentTarget.value)}
              />
            </div>

            <div class="editor-field">
              <label>Priority</label>
              <select
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
              <label>Status</label>
              <select
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
                    <span class="tag">
                      {blocker}
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
                    <span class="tag">
                      {block}
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

          <textarea
            class="editor editor-body"
            value={task()?.body || ""}
            onInput={(event) => {
              if (!task()) return
              props.onTaskChange({ ...task()!, body: event.currentTarget.value })
            }}
            placeholder="Task description and notes..."
          />
        </div>
      </Show>

      <div class="footer">
        <span>{props.status}</span>
        <span>{props.lastEvent}</span>
      </div>
    </div>
  )
}
