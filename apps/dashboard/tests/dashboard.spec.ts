import { expect, Page, test } from "@playwright/test"

type TaskItem = {
  id: string
  short_id: string
  title: string
  role: string
  priority: string
  completed: boolean
  parent: string
  blockers: string[]
  blocks: string[]
  path: string
  date_created: string
  date_edited: string
}

type TaskDetail = TaskItem & {
  body: string
}

const project = {
  name: "strandyard",
  storage_root: "/tmp/strandyard",
  tasks_root: "/tmp/strandyard/tasks",
  roles_root: "/tmp/strandyard/.strand/roles",
  templates_root: "/tmp/strandyard/templates",
  git_root: "/tmp/strandyard",
  storage: "local",
}

const tasks: TaskItem[] = [
  {
    id: "T1ab-setup-dashboard",
    short_id: "T1ab",
    title: "Set up dashboard tests",
    role: "developer",
    priority: "high",
    completed: true,
    parent: "",
    blockers: [],
    blocks: [],
    path: "tasks/T1ab-setup-dashboard/T1ab-setup-dashboard.md",
    date_created: "2026-01-20T00:00:00Z",
    date_edited: "2026-01-21T00:00:00Z",
  },
  {
    id: "T2cd-fix-stream",
    short_id: "T2cd",
    title: "Fix stream reconnect",
    role: "reviewer",
    priority: "medium",
    completed: false,
    parent: "",
    blockers: ["T1ab"],
    blocks: [],
    path: "tasks/T2cd-fix-stream/T2cd-fix-stream.md",
    date_created: "2026-01-22T00:00:00Z",
    date_edited: "2026-01-23T00:00:00Z",
  },
]

const taskDetail: TaskDetail = {
  ...tasks[0],
  body: "## Summary\nAdd coverage for the dashboard.",
}

const installEventSourceMock = async (page: Page) => {
  await page.addInitScript(() => {
    class MockEventSource {
      url: string
      onopen: ((event: Event) => void) | null = null
      onerror: ((event: Event) => void) | null = null

      constructor(url: string) {
        this.url = url
        setTimeout(() => {
          this.onopen?.(new Event("open"))
        }, 0)
      }

      addEventListener() {}

      removeEventListener() {}

      close() {}
    }

    const target = window as typeof window & { EventSource: typeof MockEventSource }
    target.EventSource = MockEventSource
  })
}

const setupApiMocks = async (
  page: Page,
  options: {
    onPatch?: (payload: Record<string, unknown>) => void
  } = {},
) => {
  await page.route("**/api/**", async (route) => {
    const request = route.request()
    const url = new URL(request.url())

    if (url.pathname === "/api/projects") {
      await route.fulfill({
        contentType: "application/json",
        body: JSON.stringify({ projects: [project], current: project.name }),
      })
      return
    }

    if (url.pathname === "/api/tasks") {
      await route.fulfill({
        contentType: "application/json",
        body: JSON.stringify(tasks),
      })
      return
    }

    if (url.pathname === "/api/task") {
      if (request.method() === "PATCH") {
        const payload = request.postDataJSON() as Record<string, unknown>
        options.onPatch?.(payload)
        const updated = {
          ...taskDetail,
          ...payload,
          date_edited: "2026-01-30T10:00:00Z",
        }
        await route.fulfill({
          contentType: "application/json",
          body: JSON.stringify(updated),
        })
        return
      }

      await route.fulfill({
        contentType: "application/json",
        body: JSON.stringify(taskDetail),
      })
      return
    }

    await route.fulfill({ status: 404, body: "Not found" })
  })
}

test("loads tasks and opens the editor", async ({ page }) => {
  await installEventSourceMock(page)
  await setupApiMocks(page)

  await page.goto("/")
  await page.getByLabel("Status").selectOption("all")

  await expect(page.getByRole("button", { name: tasks[0].title })).toBeVisible()
  await page.getByRole("button", { name: tasks[0].title }).click()

  await expect(page.getByLabel("Title")).toHaveValue(tasks[0].title)
  await expect(page.getByText(taskDetail.path)).toBeVisible()
})

test("saves edits to a task", async ({ page }) => {
  await installEventSourceMock(page)

  let patchPayload: Record<string, unknown> | undefined
  await setupApiMocks(page, {
    onPatch: (payload) => {
      patchPayload = payload
    },
  })

  await page.goto("/")
  await page.getByLabel("Status").selectOption("all")
  await page.getByRole("button", { name: tasks[0].title }).click()

  await page.getByLabel("Title").fill("Ship the dashboard tests")
  await page.getByRole("button", { name: "Save" }).click()

  await expect(page.getByText(`Saved ${tasks[0].short_id}`)).toBeVisible()
  expect(patchPayload).toMatchObject({
    title: "Ship the dashboard tests",
  })
})
