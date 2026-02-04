import { expect, test } from "@playwright/test"

const transitionPage = "/?e2eTransition=1"

test("removal triggers transitionrun and transitionend", async ({ page }) => {
  await page.goto(transitionPage)
  await page.waitForSelector('[data-testid="item-A"]')

  const removalEvents = page.evaluate(() =>
    new Promise<string[]>((resolve) => {
      const item = document.querySelector('[data-testid="item-A"]')
      if (!item) {
        resolve([])
        return
      }

      const events: string[] = []
      item.addEventListener("transitionrun", () => events.push("run"), { once: true })
      item.addEventListener(
        "transitionend",
        () => {
          events.push("end")
          resolve(events)
        },
        { once: true },
      )
    }),
  )

  await page.getByTestId("remove-first").click()

  await expect(page.getByTestId("item-A")).toHaveCount(0)
  await expect(removalEvents).resolves.toEqual(["run", "end"])
})

test("reorder triggers move transition events", async ({ page }) => {
  await page.goto(transitionPage)
  await page.waitForSelector('[data-testid="item-A"]')

  const moveEvents = page.evaluate(() =>
    new Promise<string[]>((resolve) => {
      const item = document.querySelector('[data-testid="item-A"]')
      if (!item) {
        resolve([])
        return
      }

      const events: string[] = []
      item.addEventListener("transitionrun", () => events.push("run"), { once: true })
      item.addEventListener(
        "transitionend",
        () => {
          events.push("end")
          resolve(events)
        },
        { once: true },
      )
    }),
  )

  await page.getByTestId("toggle-order").click()
  await expect(moveEvents).resolves.toEqual(["run", "end"])
})
