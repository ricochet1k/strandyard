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

const projectB: ProjectInfo = {
  name: "strandyard-alt",
  storage_root: "/tmp/strandyard-alt",
  tasks_root: "/tmp/strandyard-alt/tasks",
  roles_root: "/tmp/strandyard-alt/.strand/roles",
  templates_root: "/tmp/strandyard-alt/templates",
  git_root: "/tmp/strandyard-alt",
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

const roles = [
  { name: "developer", path: ".strand/roles/developer.md", kind: "roles" },
  { name: "reviewer", path: ".strand/roles/reviewer.md", kind: "roles" },
]

const templates = [
  { name: "task", path: "templates/task.md", kind: "templates" },
  { name: "epic", path: "templates/epic.md", kind: "templates" },
]

const roleContent = {
  path: ".strand/roles/developer.md",
  content: "---\ndescription: Implements tasks\n---\n\n# Developer Role\n\nResponsibilities:\n- Write code\n- Fix bugs",
}

const reviewerRoleContent = {
  path: ".strand/roles/reviewer.md",
  content: "---\ndescription: Reviews code\n---\n\n# Reviewer Role\n\nResponsibilities:\n- Review PRs",
}

const templateContent = {
  path: "templates/task.md",
  content: "---\nrole: developer\npriority: medium\ndescription: Basic task\nid_prefix: T\n---\n\n# Task Title",
}

const epicTemplateContent = {
  path: "templates/epic.md",
  content: "---\nrole: architect\npriority: high\ndescription: Epic template\nid_prefix: E\n---\n\n# Epic Title",
}

const selectAllShortcut = process.platform === "darwin" ? "Meta+A" : "Control+A"

const projectBTasks: TaskItem[] = [
  {
    id: "T9zz-alt-task",
    short_id: "T9zz",
    title: "Alternate project task",
    role: "ops",
    priority: "low",
    completed: false,
    parent: "",
    blockers: [],
    blocks: [],
    path: "tasks/T9zz-alt-task/T9zz-alt-task.md",
    date_created: "2026-02-01T00:00:00Z",
    date_edited: "2026-02-02T00:00:00Z",
  },
]

const projectBROpsRoleContent = {
  path: ".strand/roles/ops.md",
  content: "---\ndescription: Operates services\n---\n\n# Ops Role\n\nResponsibilities:\n- Keep services running\n- Resolve incidents",
}

const projectBTemplateContent = {
  path: "templates/ops-task.md",
  content: "---\nrole: ops\npriority: low\ndescription: Ops focused task\nid_prefix: O\n---\n\n# Ops Task",
}

const projectBRoles = [
  { name: "ops", path: projectBROpsRoleContent.path, kind: "roles" },
]

const projectBTemplates = [
  { name: "ops-task", path: projectBTemplateContent.path, kind: "templates" },
]


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
    onFilePut?: (path: string, content: string) => void
    onTaskCreate?: (payload: Record<string, unknown>) => void
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
      if (request.method() === "POST") {
        const payload = request.postDataJSON() as Record<string, unknown>
        options.onTaskCreate?.(payload)
        await route.fulfill({
          contentType: "application/json",
          status: 201,
          body: JSON.stringify({ status: "created", message: "✓ Task created: T3ef-new-task\n" }),
        })
        return
      }

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

    if (url.pathname === "/api/files") {
      const kind = url.searchParams.get("kind")
      if (kind === "roles") {
        await route.fulfill({
          contentType: "application/json",
          body: JSON.stringify(roles),
        })
        return
      }
      if (kind === "templates") {
        await route.fulfill({
          contentType: "application/json",
          body: JSON.stringify(templates),
        })
        return
      }
    }

    if (url.pathname === "/api/file") {
      const path = url.searchParams.get("path")
      
      if (request.method() === "PUT") {
        const payload = request.postDataJSON() as { content: string }
        if (path) {
          options.onFilePut?.(path, payload.content)
        }
        await route.fulfill({
          contentType: "application/json",
          body: JSON.stringify({}),
        })
        return
      }

      // Role files
      if (path === roleContent.path) {
        await route.fulfill({
          contentType: "application/json",
          body: JSON.stringify(roleContent),
        })
        return
      }

      if (path === reviewerRoleContent.path) {
        await route.fulfill({
          contentType: "application/json",
          body: JSON.stringify(reviewerRoleContent),
        })
        return
      }

      // Template files
      if (path === templateContent.path) {
        await route.fulfill({
          contentType: "application/json",
          body: JSON.stringify(templateContent),
        })
        return
      }

      if (path === epicTemplateContent.path) {
        await route.fulfill({
          contentType: "application/json",
          body: JSON.stringify(epicTemplateContent),
        })
        return
      }
    }

    await route.fulfill({ status: 404, body: "Not found" })
  })
}

const setupProjectOverrideMocks = async (page: Page) => {
  await page.route("**/api/**", async (route) => {
    const request = route.request()
    const url = new URL(request.url())
    const requestedProject = url.searchParams.get("project") || project.name

    if (url.pathname === "/api/projects") {
      await route.fulfill({
        contentType: "application/json",
        body: JSON.stringify({ projects: [project, projectB], current: project.name }),
      })
      return
    }

    if (requestedProject !== projectB.name) {
      await route.continue()
      return
    }

    if (url.pathname === "/api/tasks") {
      await route.fulfill({
        contentType: "application/json",
        body: JSON.stringify(projectBTasks),
      })
      return
    }

    if (url.pathname === "/api/files") {
      const kind = url.searchParams.get("kind")
      if (kind === "roles") {
        await route.fulfill({
          contentType: "application/json",
          body: JSON.stringify(projectBRoles),
        })
        return
      }
      if (kind === "templates") {
        await route.fulfill({
          contentType: "application/json",
          body: JSON.stringify(projectBTemplates),
        })
        return
      }
      await route.continue()
      return
    }

    if (url.pathname === "/api/file") {
      const path = url.searchParams.get("path")
      if (!path) {
        await route.continue()
        return
      }
      if (path === projectBROpsRoleContent.path) {
        await route.fulfill({
          contentType: "application/json",
          body: JSON.stringify(projectBROpsRoleContent),
        })
        return
      }
      if (path === projectBTemplateContent.path) {
        await route.fulfill({
          contentType: "application/json",
          body: JSON.stringify(projectBTemplateContent),
        })
        return
      }
      await route.continue()
      return
    }

    await route.continue()
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

test("switches to roles tab and loads role files", async ({ page }) => {
  await installEventSourceMock(page)
  await setupApiMocks(page)

  await page.goto("/")
  
  // Switch to Roles tab
  await page.getByRole("button", { name: "Roles" }).click()
  
  // Wait for roles to load
  await page.waitForTimeout(500)
  
  // Check that roles are listed in table
  await expect(page.locator(".task-table")).toBeVisible()
  await expect(page.getByRole("button", { name: "developer" })).toBeVisible()
  await expect(page.getByRole("button", { name: "reviewer" })).toBeVisible()
  
  // Click on a role to load it
  await page.getByRole("button", { name: "developer" }).click()
  
  // Check that role content is displayed in editor
  await expect(page.locator(".pane-header .detail", { hasText: roleContent.path })).toBeVisible()
  await expect(page.locator(".codemirror-container .cm-content")).toContainText("Developer Role")
  
  // Check that frontmatter controls are visible
  await expect(page.getByLabel("Description")).toHaveValue("Implements tasks")
})

test("switches to templates tab and loads template files", async ({ page }) => {
  await installEventSourceMock(page)
  await setupApiMocks(page)

  await page.goto("/")
  
  // Switch to Templates tab
  await page.getByRole("button", { name: "Templates" }).click()
  
  // Wait for templates to load
  await page.waitForTimeout(500)
  
  // Check that templates are listed in table
  await expect(page.locator(".task-table")).toBeVisible()
  await expect(page.locator(".task-table .task-link", { hasText: "task" })).toBeVisible()
  await expect(page.locator(".task-table .task-link", { hasText: "epic" })).toBeVisible()
  
  // Click on a template to load it
  await page.locator(".task-table .task-link", { hasText: "task" }).click()
  
  // Check that template content is displayed in editor
  await expect(page.locator(".pane-header .detail", { hasText: templateContent.path })).toBeVisible()
  await expect(page.locator(".codemirror-container .cm-content")).toContainText("Task Title")
  
  // Check that frontmatter controls are visible
  await expect(page.getByLabel("Role")).toHaveValue("developer")
  await expect(page.getByLabel("Priority")).toHaveValue("medium")
  await expect(page.getByLabel("Description")).toHaveValue("Basic task")
  await expect(page.getByLabel("ID Prefix")).toHaveValue("T")
})

test("saves edits to a role file", async ({ page }) => {
  await installEventSourceMock(page)

  let savedPath: string | undefined
  let savedContent: string | undefined
  await setupApiMocks(page, {
    onFilePut: (path, content) => {
      savedPath = path
      savedContent = content
    },
  })

  await page.goto("/")
  
  // Switch to Roles tab and open a role
  await page.getByRole("button", { name: "Roles" }).click()
  await page.waitForTimeout(500)
  await page.getByRole("button", { name: "developer" }).click()
  
  // Edit the role description
  await page.getByLabel("Description").fill("Updated description")
  
  // Edit the role body
  const editorContent = page.locator(".codemirror-container .cm-content")
  await editorContent.click()
  await page.keyboard.press(selectAllShortcut)
  await page.keyboard.type("# Updated Developer Role\n\nNew content")
  
  await page.getByRole("button", { name: "Save" }).click()
  
  // Check that save was called with correct data
  await expect(page.getByText(`Saved ${roleContent.path}`)).toBeVisible()
  expect(savedPath).toBe(roleContent.path)
  expect(savedContent).toContain("description: Updated description")
  expect(savedContent).toContain("# Updated Developer Role")
})

test("opens add task modal and templates populate", async ({ page }) => {
  await installEventSourceMock(page)
  
  let createdTask: Record<string, unknown> | undefined
  await setupApiMocks(page, {
    onTaskCreate: (payload) => {
      createdTask = payload
    },
  })

  await page.goto("/")
  await page.getByLabel("Status").selectOption("all")
  
  // Wait for initial load
  await page.waitForTimeout(500)
  
  // Click the Add Task button
  await page.getByRole("button", { name: "+ Add Task" }).click()
  
  // Modal should be visible
  await expect(page.getByRole("heading", { name: "Add New Task" })).toBeVisible()
  
  // Template dropdown should be visible and populated
  const templateSelect = page.locator("select#template")
  await expect(templateSelect).toBeVisible()
  
  // Check that template options are populated (check count rather than visibility)
  const templateOptions = await templateSelect.locator("option").count()
  expect(templateOptions).toBeGreaterThan(2) // At least 2 templates + placeholder
  
  // Select a template
  await templateSelect.selectOption("task")
  
  // Fill in the title
  await page.locator("input#title").fill("New test task")
  
  // Fill in description
  await page.locator("textarea#body").fill("This is a test task description")
  
  // Submit the form
  await page.getByRole("button", { name: "Create Task" }).click()
  
  // Wait for modal to close
  await expect(page.getByRole("heading", { name: "Add New Task" })).not.toBeVisible()
  
  // Verify the task was created with correct data
  expect(createdTask).toMatchObject({
    template_name: "task",
    title: "New test task",
    body: "This is a test task description",
    role: "developer", // from template default
    priority: "medium", // from template default
  })
})

test("add task modal allows selecting parent task", async ({ page }) => {
  await installEventSourceMock(page)
  
  let createdTask: Record<string, unknown> | undefined
  await setupApiMocks(page, {
    onTaskCreate: (payload) => {
      createdTask = payload
    },
  })

  await page.goto("/")
  await page.getByLabel("Status").selectOption("all")
  await page.waitForTimeout(500)
  
  // Open add task modal
  await page.getByRole("button", { name: "+ Add Task" }).click()
  
  // Select template
  await page.locator("select#template").selectOption("task")
  
  // Fill in title
  await page.locator("input#title").fill("Subtask example")
  
  // Select parent task
  const parentSelect = page.locator("select#parent")
  await expect(parentSelect).toBeVisible()
  
  // Check that parent task options are populated
  const parentOptions = await parentSelect.locator("option").count()
  expect(parentOptions).toBeGreaterThan(1) // At least 1 task + "No parent" option
  
  // Select the parent
  await parentSelect.selectOption(tasks[0].id)
  
  // Submit
  await page.getByRole("button", { name: "Create Task" }).click()
  
  // Verify parent was included
  expect(createdTask).toMatchObject({
    parent: tasks[0].id,
  })
})

test("add subtask button opens modal with parent pre-filled", async ({ page }) => {
  await installEventSourceMock(page)
  
  let createdTask: Record<string, unknown> | undefined
  await setupApiMocks(page, {
    onTaskCreate: (payload) => {
      createdTask = payload
    },
  })

  await page.goto("/")
  await page.getByLabel("Status").selectOption("all")
  await page.waitForTimeout(500)
  
  // Open a task
  await page.getByRole("button", { name: tasks[0].title }).click()
  await page.waitForTimeout(500)
  
  // Click Add Subtask button
  await page.getByRole("button", { name: "+ Add Subtask" }).click()
  
  // Modal should be visible
  await expect(page.getByRole("heading", { name: "Add New Task" })).toBeVisible()
  
  // Parent should be pre-selected
  const parentSelect = page.locator("select#parent")
  await expect(parentSelect).toHaveValue(tasks[0].id)
  
  // Fill in the form
  await page.locator("select#template").selectOption("task")
  await page.locator("input#title").fill("Subtask of first task")
  
  // Submit
  await page.getByRole("button", { name: "Create Task" }).click()
  
  // Verify the task was created with the correct parent
  expect(createdTask).toMatchObject({
    parent: tasks[0].id,
    title: "Subtask of first task",
  })
})

test("switching projects refreshes templates and roles without leaving the tasks view", async ({ page }) => {
  await installEventSourceMock(page)
  await setupApiMocks(page)
  await setupProjectOverrideMocks(page)

  await page.goto("/")
  await page.getByLabel("Status").selectOption("all")
  await page.waitForTimeout(500)

  await page.getByLabel("Project").selectOption(projectB.name)
  await page.waitForTimeout(500)

  await page.getByRole("button", { name: "+ Add Task" }).click()

  const templateSelect = page.locator("select#template")
  await expect(templateSelect.locator("option", { hasText: "ops-task" })).toHaveCount(1)

  const roleSelect = page.locator("select#role")
  await expect(roleSelect.locator("option", { hasText: "ops" })).toHaveCount(1)
})
