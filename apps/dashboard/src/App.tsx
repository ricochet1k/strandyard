import { For, Show, createEffect, createSignal, onCleanup, onMount } from "solid-js";

type Tab = "tasks" | "roles" | "templates";

type TaskItem = {
  id: string;
  short_id: string;
  title: string;
  role: string;
  priority: string;
  completed: boolean;
  parent: string;
  blockers: string[];
  blocks: string[];
  path: string;
  date_created: string;
  date_edited: string;
};

type FileEntry = {
  name: string;
  path: string;
  kind: string;
};

type FilePayload = {
  path: string;
  content: string;
};

type StreamUpdate = {
  event: string;
  path: string;
  project: string;
  task?: {
    id: string;
    file_path: string;
    title: string;
  };
};

type ProjectInfo = {
  name: string;
  storage_root: string;
  tasks_root: string;
  roles_root: string;
  templates_root: string;
  git_root: string;
  storage: string;
};

type ProjectsResponse = {
  projects: ProjectInfo[];
  current: string;
};

const tabs: { id: Tab; label: string; detail: string }[] = [
  { id: "tasks", label: "Tasks", detail: "Live task files" },
  { id: "roles", label: "Roles", detail: "Role documents" },
  { id: "templates", label: "Templates", detail: "Task templates" },
];

async function fetchJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `Request failed: ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export default function App() {
  const [tab, setTab] = createSignal<Tab>("tasks");
  const [tasks, setTasks] = createSignal<TaskItem[]>([]);
  const [files, setFiles] = createSignal<FileEntry[]>([]);
  const [activePath, setActivePath] = createSignal("");
  const [content, setContent] = createSignal("");
  const [dirty, setDirty] = createSignal(false);
  const [status, setStatus] = createSignal("");
  const [connected, setConnected] = createSignal(false);
  const [lastEvent, setLastEvent] = createSignal("");
  const [projects, setProjects] = createSignal<ProjectInfo[]>([]);
  const [currentProject, setCurrentProject] = createSignal("");

  const list = () => (tab() === "tasks" ? tasks() : files());

  const apiURL = (path: string) => {
    const project = currentProject();
    if (!project) return path;
    const separator = path.includes('?') ? '&' : '?';
    return `${path}${separator}project=${encodeURIComponent(project)}`;
  };

  const loadProjects = async () => {
    try {
      const data = await fetchJSON<ProjectsResponse>("/api/projects");
      setProjects(data.projects);
      if (!currentProject() && data.current) {
        setCurrentProject(data.current);
      }
      // Check URL params for project
      const params = new URLSearchParams(window.location.search);
      const urlProject = params.get('project');
      if (urlProject) {
        setCurrentProject(urlProject);
      }
    } catch (err) {
      setStatus(`Failed to load projects: ${errorMessage(err)}`);
    }
  };

  const loadTasks = async () => {
    try {
      const data = await fetchJSON<TaskItem[]>(apiURL("/api/tasks"));
      setTasks(data);
    } catch (err) {
      setStatus(`Failed to load tasks: ${errorMessage(err)}`);
    }
  };

  const loadFiles = async (kind: Tab) => {
    try {
      const data = await fetchJSON<FileEntry[]>(apiURL(`/api/files?kind=${kind}`));
      setFiles(data);
    } catch (err) {
      setStatus(`Failed to load ${kind}: ${errorMessage(err)}`);
    }
  };

  const loadFile = async (path: string) => {
    try {
      const data = await fetchJSON<FilePayload>(apiURL(`/api/file?path=${encodeURIComponent(path)}`));
      setActivePath(data.path);
      setContent(data.content);
      setDirty(false);
      setStatus(`Loaded ${data.path}`);
    } catch (err) {
      setStatus(`Failed to load file: ${errorMessage(err)}`);
    }
  };

  const saveFile = async () => {
    if (!activePath()) {
      return;
    }
    try {
      setStatus("Saving...");
      await fetchJSON(apiURL(`/api/file?path=${encodeURIComponent(activePath())}`), {
        method: "PUT",
        body: JSON.stringify({ content: content() }),
      });
      setDirty(false);
      setStatus(`Saved ${activePath()}`);
    } catch (err) {
      setStatus(`Save failed: ${errorMessage(err)}`);
    }
  };

  const onSelect = (entry: TaskItem | FileEntry) => {
    const path = "path" in entry ? entry.path : "";
    if (!path) {
      return;
    }
    void loadFile(path);
  };

  createEffect(() => {
    const current = tab();
    setActivePath("");
    setContent("");
    setDirty(false);
    setStatus("");
    if (current === "tasks") {
      void loadTasks();
      return;
    }
    void loadFiles(current);
  });

  createEffect(() => {
    // Reload data when project changes
    const project = currentProject();
    if (!project) return;
    setActivePath("");
    setContent("");
    setDirty(false);
    if (tab() === "tasks") {
      void loadTasks();
    } else {
      void loadFiles(tab());
    }
  });

  onMount(() => {
    void loadProjects();
    void loadTasks();

    const source = new EventSource("/api/stream");
    const onOpen = () => setConnected(true);
    const onError = () => setConnected(false);
    const onTask = (event: MessageEvent) => {
      try {
        const update = JSON.parse(event.data) as StreamUpdate;
        // Only process updates for current project
        if (update.project !== currentProject()) {
          return;
        }
        setLastEvent(`${update.event} • ${update.path}`);
        if (tab() === "tasks") {
          void loadTasks();
        }
        if (activePath() && update.path === activePath()) {
          void loadFile(activePath());
        }
      } catch (err) {
        setStatus(`Stream error: ${errorMessage(err)}`);
      }
    };

    source.addEventListener("open", onOpen);
    source.addEventListener("error", onError);
    source.addEventListener("task", onTask as EventListener);

    const keyHandler = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "s") {
        event.preventDefault();
        void saveFile();
      }
    };
    window.addEventListener("keydown", keyHandler);

    onCleanup(() => {
      source.removeEventListener("open", onOpen);
      source.removeEventListener("error", onError);
      source.removeEventListener("task", onTask as EventListener);
      source.close();
      window.removeEventListener("keydown", keyHandler);
    });
  });

  return (
    <div class="app">
      <header class="top-bar">
        <div>
          <p class="kicker">strand</p>
          <h1>Task Command Deck</h1>
          <p class="subhead">
            Inspect, edit, and collaborate on tasks with real-time file signals.
          </p>
        </div>
        <div style={{ display: "flex", gap: "1rem", "align-items": "center" }}>
          <Show when={projects().length > 1}>
            <div style={{ display: "flex", "flex-direction": "column", gap: "0.25rem" }}>
              <label style={{ "font-size": "0.75rem", "font-weight": "600" }}>Project</label>
              <select
                value={currentProject()}
                onChange={(e) => setCurrentProject(e.currentTarget.value)}
                style={{
                  padding: "0.5rem",
                  "border-radius": "4px",
                  border: "1px solid #333",
                  background: "#1a1a1a",
                  color: "#fff",
                }}
              >
                <For each={projects()}>
                  {(proj) => (
                    <option value={proj.name}>
                      {proj.name} ({proj.storage})
                    </option>
                  )}
                </For>
              </select>
            </div>
          </Show>
          <div class="status-panel">
            <div class={`status-dot ${connected() ? "online" : "offline"}`} />
            <div>
              <p class="status-title">Stream</p>
              <p class="status-detail">{connected() ? "Connected" : "Reconnecting"}</p>
            </div>
          </div>
        </div>
      </header>

      <section class="tabs">
        <For each={tabs}>
          {(item) => (
            <button
              class={`tab ${tab() === item.id ? "active" : ""}`}
              onClick={() => setTab(item.id)}
            >
              <span>{item.label}</span>
              <em>{item.detail}</em>
            </button>
          )}
        </For>
      </section>

      <section class="workspace">
        <div class="list-pane">
          <div class="pane-header">
            <h2>{tabLabel(tab())}</h2>
            <span class="pill">{list().length} items</span>
          </div>

          <div class="list">
            <For each={list()}>
              {(entry, index) => (
                <button
                  class={`list-item ${activePath() === entry.path ? "selected" : ""}`}
                  style={{ "animation-delay": `${index() * 0.03}s` }}
                  onClick={() => onSelect(entry)}
                >
                  <Show when={isTask(entry)} fallback={<RoleRow entry={entry as FileEntry} />}>
                    <TaskRow entry={entry as TaskItem} />
                  </Show>
                </button>
              )}
            </For>
          </div>
        </div>

        <div class="editor-pane">
          <div class="pane-header">
            <div>
              <h2>Editor</h2>
              <p class="detail">
                {activePath() ? activePath() : "Select a file to begin."}
              </p>
            </div>
            <div class="editor-actions">
              <span class={`sync ${dirty() ? "dirty" : ""}`}>{dirty() ? "Unsaved" : "Synced"}</span>
              <button class="primary" disabled={!activePath()} onClick={() => void saveFile()}>
                Save
              </button>
            </div>
          </div>

          <textarea
            class="editor"
            value={content()}
            onInput={(event) => {
              setContent(event.currentTarget.value);
              setDirty(true);
            }}
            placeholder="Select a task, role, or template file to edit."
          />

          <div class="footer">
            <span>{status()}</span>
            <span>{lastEvent()}</span>
          </div>
        </div>
      </section>
    </div>
  );
}

function TaskRow(props: { entry: TaskItem }) {
  const entry = () => props.entry;
  return (
    <div class="task-row">
      <div>
        <p class="task-title">{entry().title || "Untitled task"}</p>
        <p class="task-meta">
          <span>{entry().short_id}</span>
          <span>{entry().role || "unassigned"}</span>
          <span>{entry().completed ? "done" : "active"}</span>
        </p>
      </div>
      <span class={`priority ${entry().priority}`}>{entry().priority}</span>
    </div>
  );
}

function RoleRow(props: { entry: FileEntry }) {
  const entry = () => props.entry;
  return (
    <div class="role-row">
      <div>
        <p class="task-title">{entry().name}</p>
        <p class="task-meta">
          <span>{entry().kind}</span>
          <span>{entry().path}</span>
        </p>
      </div>
    </div>
  );
}

function tabLabel(tab: Tab) {
  switch (tab) {
    case "tasks":
      return "Tasks Library";
    case "roles":
      return "Roles Library";
    case "templates":
      return "Templates Library";
    default:
      return "Library";
  }
}

function isTask(entry: TaskItem | FileEntry): entry is TaskItem {
  return "short_id" in entry;
}

function errorMessage(err: unknown) {
  if (err instanceof Error) {
    return err.message;
  }
  return "Unknown error";
}
