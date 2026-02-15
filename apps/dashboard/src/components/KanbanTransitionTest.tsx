import { createMemo, createSignal, For } from "solid-js"
import { MyTransitionGroup } from "./MyTransitionGroup"
import "./KanbanBoard.css"

type Item = {
  id: string
  status: "open" | "done"
}

const columns: Array<Item["status"]> = ["open", "done"]

export default function KanbanTransitionTest() {
  const [items, setItems] = createSignal<Item[]>([
    { id: "A", status: "open" },
    { id: "B", status: "open" },
    { id: "C", status: "done" },
  ])

  const grouped = createMemo(() => {
    const map: Record<Item["status"], Item[]> = { open: [], done: [] }
    for (const item of items()) map[item.status].push(item)
    return map
  })

  const backendMove = () => {
    setItems((current) =>
      current.map((item) =>
        item.id === "A" ? { ...item, status: item.status === "open" ? "done" : "open" } : item,
      ),
    )
  }

  return (
    <div style={{ padding: "24px" }}>
      <button data-testid="backend-move" onClick={backendMove}>
        Backend move A
      </button>
      <div class="kanban-board" style={{ "margin-top": "12px", "grid-template-columns": "repeat(2, minmax(180px, 1fr))" }}>
        <For each={columns}>
          {(column) => (
            <section class="kanban-column" data-testid={`test-column-${column}`}>
              <header class="kanban-column-header">
                <h3>{column}</h3>
                <span>{grouped()[column].length}</span>
              </header>
              <div class="kanban-column-body">
                <MyTransitionGroup classPrefix="kanban-card">
                  <For each={grouped()[column]}>
                    {(item) => (
                      <article class="kanban-card" data-testid={`test-card-${item.id}`}>
                        {item.id}
                      </article>
                    )}
                  </For>
                </MyTransitionGroup>
              </div>
            </section>
          )}
        </For>
      </div>
    </div>
  )
}
