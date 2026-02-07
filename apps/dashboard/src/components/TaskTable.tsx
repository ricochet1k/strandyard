import { For, Show } from "solid-js"
import "./TaskTable.css"
import { MyTransitionGroup } from "./MyTransitionGroup"
import { SortDirection, SortField, TaskTreeNode } from "../App"

type TaskTableProps = {
  tasks: TaskTreeNode[]
  activePath: string
  sortField: SortField
  sortDirection: SortDirection
  viewMode: "tree" | "list"
  hasChildren: (node: TaskTreeNode) => boolean
  onSelect: (node: TaskTreeNode) => void
  isExpanded: (node: TaskTreeNode) => boolean
  onToggleNode: (node: TaskTreeNode) => void
  onSortChange: (field: SortField) => void
}

export default function TaskTable(props: TaskTableProps) {
  const handleSort = (field: SortField) => {
    props.onSortChange(field)
  }

  return (
    <table class="task-table" style="height: 100%;">
      <thead>
        <tr>
          <th class="task-col-title">
            <button
              class={`task-col-header ${props.sortField === "title" ? "active" : ""}`}
              onClick={() => handleSort("title")}
            >
              Title {props.sortField === "title" && (props.sortDirection === "asc" ? "↑" : "↓")}
            </button>
          </th>
          <th class="task-col-id">
            <button
              class={`task-col-header ${props.sortField === "short_id" ? "active" : ""}`}
              onClick={() => handleSort("short_id")}
            >
              ID {props.sortField === "short_id" && (props.sortDirection === "asc" ? "↑" : "↓")}
            </button>
          </th>
          <th class="task-col-role">
            <button
              class={`task-col-header ${props.sortField === "date_edited" ? "active" : ""}`}
              onClick={() => handleSort("date_edited")}
            >
              Role {props.sortField === "date_edited" && (props.sortDirection === "asc" ? "↑" : "↓")}
            </button>
          </th>
          <th class="task-col-status">
            <button
              class={`task-col-header ${props.sortField === "priority" ? "active" : ""}`}
              onClick={() => handleSort("priority")}
            >
              Status {props.sortField === "priority" && (props.sortDirection === "asc" ? "↑" : "↓")}
            </button>
          </th>
          <th class="task-col-priority">
            <button
              class={`task-col-header ${props.sortField === "priority" ? "active" : ""}`}
              onClick={() => handleSort("priority")}
            >
              Priority {props.sortField === "priority" && (props.sortDirection === "asc" ? "↑" : "↓")}
            </button>
          </th>
        </tr>
      </thead>
      <tbody style="overflow: auto; height: 100%;">
        <MyTransitionGroup classPrefix="tree-item">
          <For each={props.tasks}>
            {(node) => (
              <tr
                class={`tree-item ${props.activePath === node.task.path ? "selected" : ""} ${node.task.completed ? "done" : ""
                  } ${node.task.blockers && node.task.blockers.length > 0 ? "blocked" : ""}`}
              >
                <td class="task-col-title" style={{ "padding-left": props.viewMode === "tree" ? `${node.depth * 12}px` : "0" }}>
                  {(() => { console.log("task rendering..."); return undefined })()}
                  <div class="task-title-cell">
                    <Show when={props.viewMode === "tree"}>
                      <div class="task-toggle-spacer">
                        <Show when={props.hasChildren(node)}>
                          <button
                            class="tree-toggle"
                            onClick={(e) => {
                              e.stopPropagation()
                              props.onToggleNode(node)
                            }}
                          >
                            {props.isExpanded(node) ? "−" : "+"}
                          </button>
                        </Show>
                      </div>
                    </Show>
                    <button
                      class="task-link"
                      title={node.task.title}
                      onClick={() => props.onSelect(node)}
                    >
                      {node.task.title || "Untitled task"}
                    </button>
                  </div>
                </td>
                <td class="task-col-id">
                  <span>{node.task.short_id}</span>
                </td>
                <td class="task-col-role">
                  <span>{node.task.role || "unassigned"}</span>
                </td>
                <td class="task-col-status">
                  <span
                    class={`status-badge ${node.task.completed ? "status-done" : node.task.blockers?.length ? "status-blocked" : "status-active"
                      }`}
                  >
                    {node.task.completed ? "done" : node.task.blockers?.length ? "blocked" : "active"}
                  </span>
                </td>
                <td class="task-col-priority">
                  <span class={`priority ${node.task.priority}`}>{node.task.priority}</span>
                </td>
              </tr>
            )}
          </For>
        </MyTransitionGroup>
      </tbody>
    </table>
  )
}
