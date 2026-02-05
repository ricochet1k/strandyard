import { createSignal, For, Show } from "solid-js"
import type { TemplateItem, RoleItem } from "../App"
import "./AddTaskModal.css"

type AddTaskModalProps = {
  templates: TemplateItem[]
  roles: RoleItem[]
  tasks: Array<{ id: string; short_id: string; title: string }>
  defaultParent?: string
  onClose: () => void
  onSubmit: (data: {
    template_name: string
    title: string
    role: string
    priority: string
    parent: string
    body: string
  }) => Promise<void>
}

export default function AddTaskModal(props: AddTaskModalProps) {
  const [templateName, setTemplateName] = createSignal("")
  const [title, setTitle] = createSignal("")
  const [role, setRole] = createSignal("")
  const [priority, setPriority] = createSignal("medium")
  const [parent, setParent] = createSignal(props.defaultParent || "")
  const [body, setBody] = createSignal("")
  const [submitting, setSubmitting] = createSignal(false)
  const [error, setError] = createSignal("")

  const selectedTemplate = () => {
    const name = templateName()
    return props.templates.find((t) => t.name === name)
  }

  const effectiveRole = () => {
    const r = role()
    if (r) return r
    const template = selectedTemplate()
    return template?.role || ""
  }

  const effectivePriority = () => {
    const p = priority()
    if (p) return p
    const template = selectedTemplate()
    return template?.priority || "medium"
  }

  const handleSubmit = async (e: Event) => {
    e.preventDefault()
    setError("")

    const tmpl = templateName()
    const t = title().trim()

    if (!tmpl) {
      setError("Please select a template")
      return
    }
    if (!t) {
      setError("Please provide a title")
      return
    }

    setSubmitting(true)
    try {
      await props.onSubmit({
        template_name: tmpl,
        title: t,
        role: effectiveRole(),
        priority: effectivePriority(),
        parent: parent(),
        body: body(),
      })
      props.onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create task")
      setSubmitting(false)
    }
  }

  const handleBackdropClick = (e: MouseEvent) => {
    if (e.target === e.currentTarget) {
      props.onClose()
    }
  }

  return (
    <div class="modal-backdrop" onClick={handleBackdropClick}>
      <div class="modal-dialog">
        <div class="modal-header">
          <h2>Add New Task</h2>
          <button class="modal-close" onClick={props.onClose}>
            ×
          </button>
        </div>

        <form onSubmit={handleSubmit}>
          <div class="modal-body">
            <Show when={error()}>
              <div class="error-message">{error()}</div>
            </Show>

            <div class="form-group">
              <label for="template">Template *</label>
              <select
                id="template"
                value={templateName()}
                onInput={(e) => setTemplateName(e.currentTarget.value)}
                required
              >
                <option value="">-- Select a template --</option>
                <For each={props.templates}>
                  {(template) => (
                    <option value={template.name}>
                      {template.name}
                      {template.description ? ` - ${template.description}` : ""}
                    </option>
                  )}
                </For>
              </select>
            </div>

            <div class="form-group">
              <label for="title">Title *</label>
              <input
                id="title"
                type="text"
                value={title()}
                onInput={(e) => setTitle(e.currentTarget.value)}
                placeholder="Enter task title"
                required
              />
            </div>

            <div class="form-group">
              <label for="role">
                Role {selectedTemplate()?.role && `(default: ${selectedTemplate()!.role})`}
              </label>
              <select
                id="role"
                value={role()}
                onInput={(e) => setRole(e.currentTarget.value)}
              >
                <option value="">-- Use template default --</option>
                <For each={props.roles}>
                  {(r) => (
                    <option value={r.name}>
                      {r.name}
                      {r.description ? ` - ${r.description}` : ""}
                    </option>
                  )}
                </For>
              </select>
            </div>

            <div class="form-group">
              <label for="priority">
                Priority {selectedTemplate()?.priority && `(default: ${selectedTemplate()!.priority})`}
              </label>
              <select
                id="priority"
                value={priority()}
                onInput={(e) => setPriority(e.currentTarget.value)}
              >
                <option value="high">High</option>
                <option value="medium">Medium</option>
                <option value="low">Low</option>
              </select>
            </div>

            <div class="form-group">
              <label for="parent">Parent Task (optional)</label>
              <select
                id="parent"
                value={parent()}
                onInput={(e) => setParent(e.currentTarget.value)}
              >
                <option value="">-- No parent --</option>
                <For each={props.tasks}>
                  {(task) => (
                    <option value={task.id}>
                      {task.short_id} - {task.title}
                    </option>
                  )}
                </For>
              </select>
            </div>

            <div class="form-group">
              <label for="body">Description (optional)</label>
              <textarea
                id="body"
                value={body()}
                onInput={(e) => setBody(e.currentTarget.value)}
                placeholder="Additional details or context"
                rows={6}
              />
            </div>
          </div>

          <div class="modal-footer">
            <button
              type="button"
              class="button button-secondary"
              onClick={props.onClose}
              disabled={submitting()}
            >
              Cancel
            </button>
            <button
              type="submit"
              class="button button-primary"
              disabled={submitting()}
            >
              {submitting() ? "Creating..." : "Create Task"}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
