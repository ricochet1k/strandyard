import { For, Show } from "solid-js"
import "./Sidebar.css"

type Tab = "tasks" | "roles" | "templates"

type SidebarProps = {
  tab: Tab
  projects: any[]
  currentProject: string
  searchQuery: string
  filterStatus: "all" | "active" | "done"
  filterRole: string
  filterPriority: string
  availableRoles: string[]
  availablePriorities: string[]
  onTabChange: (tab: Tab) => void
  onProjectChange: (project: string) => void
  onSearchChange: (query: string) => void
  onFilterStatusChange: (status: "all" | "active" | "done") => void
  onFilterRoleChange: (role: string) => void
  onFilterPriorityChange: (priority: string) => void
}

const tabs: { id: Tab; label: string; detail: string }[] = [
  { id: "tasks", label: "Tasks", detail: "Live task files" },
  { id: "roles", label: "Roles", detail: "Role documents" },
  { id: "templates", label: "Templates", detail: "Task templates" },
]

export default function Sidebar(props: SidebarProps) {
  return (
    <aside class="sidebar">
      {/* Project Selector */}
      <Show when={props.projects.length > 1}>
        <div class="sidebar-filter">
          <label>Project</label>
          <select
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
          <p class="sidebar-section-title">Filters</p>
          <div class="sidebar-filter">
            <label for="filter-status">Status</label>
            <select
              id="filter-status"
              value={props.filterStatus}
              onChange={(e) => props.onFilterStatusChange(e.currentTarget.value as any)}
            >
              <option value="all">All Status</option>
              <option value="active">Active</option>
              <option value="done">Done</option>
            </select>
          </div>

          <div class="sidebar-filter">
            <label for="filter-role">Role</label>
            <select
              id="filter-role"
              value={props.filterRole}
              onChange={(e) => props.onFilterRoleChange(e.currentTarget.value)}
            >
              <option value="all">All Roles</option>
              <For each={props.availableRoles}>
                {(role) => <option value={role}>{role}</option>}
              </For>
            </select>
          </div>

          <div class="sidebar-filter">
            <label for="filter-priority">Priority</label>
            <select
              id="filter-priority"
              value={props.filterPriority}
              onChange={(e) => props.onFilterPriorityChange(e.currentTarget.value)}
            >
              <option value="all">All Priorities</option>
              <For each={props.availablePriorities}>
                {(priority) => <option value={priority}>{priority}</option>}
              </For>
            </select>
          </div>
        </div>
      </Show>
    </aside>
  )
}
