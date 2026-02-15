import { combine } from "@atlaskit/pragmatic-drag-and-drop/combine"
import { draggable, dropTargetForElements } from "@atlaskit/pragmatic-drag-and-drop/element/adapter"
import { createMemo, createSignal, For, onCleanup } from "solid-js"
import { TaskTreeNode } from "../App"
import { MyTransitionGroup } from "./MyTransitionGroup"
import "./KanbanBoard.css"

type KanbanStatus = "open" | "in_progress" | "done" | "cancelled" | "duplicate"

type KanbanBoardProps = {
  tasks: TaskTreeNode[]
  activePath: string
  onSelect: (node: TaskTreeNode) => void
  onMoveTask: (taskID: string, nextStatus: KanbanStatus) => Promise<void>
}

type CardData = {
  type: "kanban-card"
  taskID: string
  fromStatus: KanbanStatus
}

const columns: { id: KanbanStatus; label: string }[] = [
  { id: "open", label: "Open" },
  { id: "in_progress", label: "In Progress" },
  { id: "done", label: "Done" },
  { id: "cancelled", label: "Cancelled" },
  { id: "duplicate", label: "Duplicate" },
]

function normalizeTaskStatus(task: { status: string; completed: boolean }): KanbanStatus {
  const status = task.status?.trim() as KanbanStatus
  if (status && columns.some((column) => column.id === status)) return status
  return task.completed ? "done" : "open"
}

function isCardData(value: Record<string, unknown>): value is CardData {
  return value.type === "kanban-card" && typeof value.taskID === "string" && typeof value.fromStatus === "string"
}

export default function KanbanBoard(props: KanbanBoardProps) {
  const [draggedTaskID, setDraggedTaskID] = createSignal<string | null>(null)
  const [hoverColumn, setHoverColumn] = createSignal<KanbanStatus | null>(null)
  const [pendingTaskID, setPendingTaskID] = createSignal<string | null>(null)
  const cleanups: Array<() => void> = []

  const grouped = createMemo(() => {
    const byColumn: Record<KanbanStatus, TaskTreeNode[]> = {
      open: [],
      in_progress: [],
      done: [],
      cancelled: [],
      duplicate: [],
    }
    for (const node of props.tasks) {
      byColumn[normalizeTaskStatus(node.task)].push(node)
    }
    return byColumn
  })

  const moveTask = (taskID: string, fromStatus: KanbanStatus, columnStatus: KanbanStatus) => {
    if (fromStatus === columnStatus) return
    if (pendingTaskID() === taskID) return
    setPendingTaskID(taskID)
    void props.onMoveTask(taskID, columnStatus).finally(() => {
      setPendingTaskID((current) => (current === taskID ? null : current))
    })
  }

  const wireColumnDropTarget = (element: HTMLElement, columnStatus: KanbanStatus) => {
    const onNativeDragOver = (event: DragEvent) => {
      event.preventDefault()
    }
    const onNativeDrop = (event: DragEvent) => {
      event.preventDefault()
      const taskID = event.dataTransfer?.getData("text/kanban-task-id")
      const fromStatus = event.dataTransfer?.getData("text/kanban-task-status") as KanbanStatus
      if (!taskID || !fromStatus) return
      setHoverColumn(null)
      moveTask(taskID, fromStatus, columnStatus)
    }
    element.addEventListener("dragover", onNativeDragOver)
    element.addEventListener("drop", onNativeDrop)

    cleanups.push(
      dropTargetForElements({
        element,
        canDrop: ({ source }) => {
          const data = source.data
          if (!isCardData(data)) return false
          return data.fromStatus !== columnStatus
        },
        onDragEnter: () => setHoverColumn(columnStatus),
        onDragLeave: () => setHoverColumn((current) => (current === columnStatus ? null : current)),
        onDrop: ({ source }) => {
          setHoverColumn(null)
          const data = source.data
          if (!isCardData(data) || data.fromStatus === columnStatus) return
          moveTask(data.taskID, data.fromStatus, columnStatus)
        },
      }),
    )
    cleanups.push(() => {
      element.removeEventListener("dragover", onNativeDragOver)
      element.removeEventListener("drop", onNativeDrop)
    })
  }

  const wireCardDraggable = (element: HTMLElement, node: TaskTreeNode) => {
    const fromStatus = normalizeTaskStatus(node.task)
    const onNativeDragStart = (event: DragEvent) => {
      setDraggedTaskID(node.task.id)
      if (!event.dataTransfer) return
      event.dataTransfer.effectAllowed = "move"
      event.dataTransfer.setData("text/kanban-task-id", node.task.id)
      event.dataTransfer.setData("text/kanban-task-status", fromStatus)
    }
    const onNativeDragEnd = () => {
      setDraggedTaskID(null)
    }
    element.setAttribute("draggable", "true")
    element.addEventListener("dragstart", onNativeDragStart)
    element.addEventListener("dragend", onNativeDragEnd)

    cleanups.push(
      combine(
        draggable({
          element,
          getInitialData: () => ({ type: "kanban-card", taskID: node.task.id, fromStatus }),
          onDragStart: () => setDraggedTaskID(node.task.id),
          onDrop: () => setDraggedTaskID(null),
        }),
      ),
    )
    cleanups.push(() => {
      element.removeEventListener("dragstart", onNativeDragStart)
      element.removeEventListener("dragend", onNativeDragEnd)
    })
  }

  onCleanup(() => {
    for (const cleanup of cleanups) cleanup()
  })

  return (
    <div class="kanban-board" data-testid="kanban-board">
      <For each={columns}>
        {(column) => (
          <section
            class={`kanban-column ${hoverColumn() === column.id ? "is-drop-target" : ""}`}
            data-testid={`kanban-column-${column.id}`}
            ref={(el) => wireColumnDropTarget(el, column.id)}
          >
            <header class="kanban-column-header">
              <h3>{column.label}</h3>
              <span>{grouped()[column.id].length}</span>
            </header>
            <div class="kanban-column-body">
              <MyTransitionGroup classPrefix="kanban-card">
                <For each={grouped()[column.id]}>
                  {(node) => (
                    <article
                      class={`kanban-card ${props.activePath === node.task.path ? "selected" : ""} ${draggedTaskID() === node.task.id ? "is-dragging" : ""} ${pendingTaskID() === node.task.id ? "is-pending" : ""}`}
                      data-testid={`kanban-card-${node.task.short_id}`}
                      ref={(el) => wireCardDraggable(el, node)}
                    >
                      <button class="kanban-card-title" title={node.task.title} onClick={() => props.onSelect(node)}>
                        {node.task.title}
                      </button>
                      <div class="kanban-card-meta">
                        <span>{node.task.short_id}</span>
                        <span>{node.task.role || "unassigned"}</span>
                        <span class={`priority ${node.task.priority}`}>{node.task.priority || "none"}</span>
                      </div>
                    </article>
                  )}
                </For>
              </MyTransitionGroup>
            </div>
          </section>
        )}
      </For>
    </div>
  )
}
