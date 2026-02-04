import { Show } from "solid-js"
import "./Editor.css"

type ParsedTask = {
  frontmatter: {
    title?: string
    role?: string
    priority?: string
    completed?: boolean
    [key: string]: any
  }
  body: string
}

type EditorProps = {
  activePath: string
  content: string
  dirty: boolean
  status: string
  lastEvent: string
  tab: string
  parsedTask: ParsedTask
  onContentChange: (content: string) => void
  onParsedTaskChange: (task: ParsedTask) => void
  onSave: () => void
}

export default function Editor(props: EditorProps) {
  const updateFrontmatterField = (field: string, value: any) => {
    const parsed = { ...props.parsedTask }
    parsed.frontmatter = { ...parsed.frontmatter, [field]: value }
    props.onParsedTaskChange(parsed)
  }

  return (
    <div class="editor-pane">
      <div class="pane-header">
        <div>
          <h2>Editor</h2>
          <p class="detail">
            {props.activePath ? props.activePath : "Select a file to begin."}
          </p>
        </div>
        <div class="editor-actions">
          <span class={`sync ${props.dirty ? "dirty" : ""}`}>
            {props.dirty ? "Unsaved" : "Synced"}
          </span>
          <button class="primary" disabled={!props.activePath} onClick={() => props.onSave()}>
            Save
          </button>
        </div>
      </div>

      <Show
        when={props.tab === "tasks" && props.activePath}
        fallback={
          <textarea
            class="editor"
            value={props.content}
            onInput={(event) => props.onContentChange(event.currentTarget.value)}
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
                value={props.parsedTask.frontmatter.title || ""}
                onInput={(e) => updateFrontmatterField("title", e.currentTarget.value)}
              />
            </div>

            <div class="editor-field">
              <label>Role</label>
              <input
                type="text"
                class="editor-input"
                value={props.parsedTask.frontmatter.role || ""}
                onInput={(e) => updateFrontmatterField("role", e.currentTarget.value)}
              />
            </div>

            <div class="editor-field">
              <label>Priority</label>
              <select
                class="editor-input"
                value={props.parsedTask.frontmatter.priority || ""}
                onChange={(e) => updateFrontmatterField("priority", e.currentTarget.value)}
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
                value={props.parsedTask.frontmatter.completed ? "done" : "active"}
                onChange={(e) => updateFrontmatterField("completed", e.currentTarget.value === "done")}
              >
                <option value="active">Active</option>
                <option value="done">Done</option>
              </select>
            </div>
          </div>

          <textarea
            class="editor editor-body"
            value={props.parsedTask.body}
            onInput={(event) => {
              const parsed = { ...props.parsedTask }
              parsed.body = event.currentTarget.value
              props.onParsedTaskChange(parsed)
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
