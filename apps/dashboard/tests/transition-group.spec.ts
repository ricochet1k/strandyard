import { expect, test } from "@playwright/test"

const transitionPage = "/?e2eTransition=1"
const kanbanTransitionPage = "/?e2eKanban=1"

test("removal emits transitionend before cleanup", async ({ page }) => {
  await page.goto(transitionPage)
  await page.waitForSelector('[data-testid="item-A"]')

  const removalEnded = page.evaluate(() =>
    new Promise<string>((resolve) => {
      const item = document.querySelector('[data-testid="item-A"]')
      if (!item) {
        resolve("missing")
        return
      }

      item.addEventListener(
        "transitionend",
        () => {
          resolve("ended")
        },
        { once: true },
      )
    }),
  )

  await page.getByTestId("remove-first").click()

  await expect(page.getByTestId("item-A")).toHaveCount(0)
  await expect(removalEnded).resolves.toBe("ended")
})

test("reorder emits transitionend for moved item", async ({ page }) => {
  await page.goto(transitionPage)
  await page.waitForSelector('[data-testid="item-A"]')

  const moveEnded = page.evaluate(() =>
    new Promise<string>((resolve) => {
      const item = document.querySelector('[data-testid="item-A"]')
      if (!item) {
        resolve("missing")
        return
      }

      item.addEventListener(
        "transitionend",
        () => {
          resolve("ended")
        },
        { once: true },
      )
    }),
  )

  await page.getByTestId("toggle-order").click()
  await expect(moveEnded).resolves.toBe("ended")
})

test("kanban backend move emits transitionend on exit", async ({ page }) => {
  await page.goto(kanbanTransitionPage)
  await page.waitForSelector('[data-testid="test-card-A"]')

  const exitEnded = page.evaluate(() =>
    new Promise<string>((resolve) => {
      const item = document.querySelector('[data-testid="test-column-open"] [data-testid="test-card-A"]')
      if (!item) {
        resolve("missing")
        return
      }
      item.addEventListener(
        "transitionend",
        () => {
          resolve("ended")
        },
        { once: true },
      )
    }),
  )

  await page.getByTestId("backend-move").click()
  await expect(page.locator('[data-testid="test-column-done"] [data-testid="test-card-A"]')).toHaveCount(1)
  await expect(exitEnded).resolves.toBe("ended")
})
