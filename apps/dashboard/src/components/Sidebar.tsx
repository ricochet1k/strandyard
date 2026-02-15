import { For, Show } from "solid-js"
import "./Sidebar.css"

type Tab = "tasks" | "roles" | "templates"

type SidebarProps = {
  tab: Tab
  projects: any[]
  currentProject: string
  searchQuery: string
  filterStatus: string[]
  filterRole: string[]
  filterPriority: string[]
  hideBlocked: boolean
  viewMode: "tree" | "list" | "kanban"
  availableRoles: string[]
  availablePriorities: string[]
  onTabChange: (tab: Tab) => void
  onProjectChange: (project: string) => void
  onSearchChange: (query: string) => void
  onFilterStatusChange: (status: string[]) => void
  onFilterRoleChange: (role: string[]) => void
  onFilterPriorityChange: (priority: string[]) => void
  onHideBlockedChange: (hide: boolean) => void
  onViewModeChange: (mode: "tree" | "list" | "kanban") => void
}

const tabs: { id: Tab; label: string; detail: string }[] = [
  { id: "tasks", label: "Tasks", detail: "Live task files" },
  { id: "roles", label: "Roles", detail: "Role documents" },
  { id: "templates", label: "Templates", detail: "Task templates" },
]

export default function Sidebar(props: SidebarProps) {
  const statusOptions = ["active", "open", "in_progress", "done", "cancelled", "duplicate"]

  const toggleFilterValue = (current: string[], value: string, checked: boolean) => {
    if (checked) {
      if (current.includes(value)) return current
      return [...current, value]
    }
    return current.filter((item) => item !== value)
  }

  const renderMultiSelect = (options: string[], selected: string[], onChange: (next: string[]) => void, label: string) => {
    const allSelected = selected.length === 0
    return (
      <div class="filter-group">
        <label class="filter-group-label">{label}</label>
        <label class="filter-option">
          <input
            type="checkbox"
            checked={allSelected}
            onChange={(e) => {
              if (e.currentTarget.checked) onChange([])
            }}
          />
          <span>All {label}</span>
        </label>
        <div class="filter-options">
          <For each={options}>
            {(option) => (
              <label class="filter-option">
                <input
                  type="checkbox"
                  checked={!allSelected && selected.includes(option)}
                  onChange={(e) => {
                    const next = toggleFilterValue(allSelected ? [] : selected, option, e.currentTarget.checked)
                    onChange(next)
                  }}
                />
                <span>{option.replace(/_/g, " ")}</span>
              </label>
            )}
          </For>
        </div>
      </div>
    )
  }

  return (
    <aside class="sidebar">
      {/* Project Selector */}
      <Show when={props.projects.length > 1}>
        <div class="sidebar-filter">
        <label for="project-select">Project</label>
        <select
          id="project-select"
          value={props.currentProject}
          onChange={(e) => props.onProjectChange(e.currentTarget.value)}
        >
            <For each={props.projects}>
              {(proj) => (
                <option value={proj.name}>{proj.name}</option>
              )}
            </For>
          </select>
        </div>
      </Show>

      {/* Navigation */}
      <div class="sidebar-section">
        <p class="sidebar-section-title">Navigate</p>
        <div class="sidebar-nav">
          <For each={tabs}>
            {(item) => (
              <button
                class={`sidebar-nav-item ${props.tab === item.id ? "active" : ""}`}
                onClick={() => props.onTabChange(item.id)}
              >
                {item.label}
              </button>
            )}
          </For>
        </div>
      </div>

      {/* Filters */}
      <Show when={props.tab === "tasks"}>
        <div class="sidebar-section">
          <p class="sidebar-section-title">Search</p>
          <input
            type="text"
            class="search-input"
            placeholder="Search tasks..."
            value={props.searchQuery}
            onInput={(e) => props.onSearchChange(e.currentTarget.value)}
          />
        </div>

        <div class="sidebar-section">
          <p class="sidebar-section-title">View</p>
          <div class="sidebar-filter">
            <label for="view-mode">Mode</label>
            <select
              id="view-mode"
              value={props.viewMode}
              onChange={(e) => props.onViewModeChange(e.currentTarget.value as any)}
            >
              <option value="tree">Tree</option>
              <option value="list">List</option>
              <option value="kanban">Kanban</option>
            </select>
          </div>
        </div>

        <div class="sidebar-section">
          <p class="sidebar-section-title">Filters</p>
          {renderMultiSelect(statusOptions, props.filterStatus, props.onFilterStatusChange, "Status")}
          {renderMultiSelect(props.availableRoles, props.filterRole, props.onFilterRoleChange, "Roles")}
          {renderMultiSelect(props.availablePriorities, props.filterPriority, props.onFilterPriorityChange, "Priorities")}

          <div class="sidebar-filter">
            <label>
              <input
                type="checkbox"
                checked={props.hideBlocked}
                onChange={(e) => props.onHideBlockedChange(e.currentTarget.checked)}
              />
              <span style={{ "margin-left": "0.5rem" }}>Hide blocked tasks</span>
            </label>
          </div>
        </div>
      </Show>
    </aside>
  )
}
