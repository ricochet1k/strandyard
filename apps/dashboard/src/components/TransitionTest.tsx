import { createSignal, For } from "solid-js"
import "./TransitionTest.css"
import { MyTransitionGroup } from "./MyTransitionGroup"

export default function TransitionTest() {
  const [items, setItems] = createSignal(["A", "B", "C"])

  const removeFirst = () => {
    setItems((current) => current.slice(1))
  }

  const toggleOrder = () => {
    setItems((current) => [...current].reverse())
  }

  const replaceAll = () => {
    setItems(Array.from(Array(5)).map(_ => "ABCDEFGHIJ"[(Math.random() * 10) | 0]))
  }

  return (
    <div class="transition-test">
      <div class="transition-test-actions">
        <button data-testid="remove-first" onClick={removeFirst}>Remove First</button>
        <button data-testid="toggle-order" onClick={toggleOrder}>Toggle Order</button>
        <button data-testid="replace-all" onClick={replaceAll}>Replace All</button>
      </div>
      <div class="transition-test-list">
        <MyTransitionGroup classPrefix="transition-test">
          <For each={items()}>
            {(item) => (
              <div class="transition-test-item" data-testid={`item-${item}`}>
                {item}
              </div>
            )}
          </For>
        </MyTransitionGroup>
      </div>
    </div>
  )
}
