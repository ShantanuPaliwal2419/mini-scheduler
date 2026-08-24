
# Mini Scheduler

A workflow scheduler built from scratch in Go to understand how modern
workflow engines execute dependency-driven tasks, prioritize work,
coordinate concurrent execution, and eventually distribute workloads across
multiple workers.

The project is intentionally built incrementally.

The goal is not to build a production replacement for Airflow or Temporal,
but to understand the core engineering mechanisms behind workflow
orchestration.

---

# Why Build This?

Workflow engines fundamentally solve a dependency and scheduling problem.

Given:

```text
        A
       / \
      B   C
       \ /
        D
```

the scheduler needs to determine:

-   Which tasks are ready?
    
-   Which tasks are blocked?
    
-   Which ready task should execute first?
    
-   Which tasks can execute concurrently?
    
-   What happens when a task fails?
    
-   What happens when a worker crashes?
    
-   How can another worker safely continue the work?
    

This project builds those mechanisms from first principles.

----------

# Final Architecture

The final goal is to evolve the scheduler into a concurrent and distributed  
workflow execution engine.

```text
                              ┌──────────────────────┐
                              │     Workflow API     │
                              │                      │
                              │ Submit Workflow      │
                              │ Inspect Workflow     │
                              │ Query Task State     │
                              └──────────┬───────────┘
                                         │
                                         ▼
                              ┌──────────────────────┐
                              │   Workflow Manager   │
                              │                      │
                              │ DAG                  │
                              │ Task State           │
                              │ Workflow Lifecycle   │
                              └──────────┬───────────┘
                                         │
                         ┌───────────────┴───────────────┐
                         │                               │
                         ▼                               ▼
              ┌──────────────────────┐       ┌──────────────────────┐
              │   Dependency Engine  │       │  Ready Task Detector │
              │                      │       │                      │
              │ Parents              │       │ Dependency Checks    │
              │ Children             │       │ State Checks         │
              │ DAG Validation       │       │ Ready Tasks          │
              │ Cycle Detection      │       └──────────┬───────────┘
              └──────────────────────┘                  │
                                                        ▼
                                             ┌──────────────────────┐
                                             │  Priority Scheduler  │
                                             │                      │
                                             │ Priority             │
                                             │ SLA                  │
                                             │ Scheduling Policy    │
                                             └──────────┬───────────┘
                                                        │
                                                        ▼
                                             ┌──────────────────────┐
                                             │     Task Queue       │
                                             │                      │
                                             │ Ready Tasks          │
                                             │ Priority Ordering    │
                                             └──────────┬───────────┘
                                                        │
                              ┌─────────────────────────┼─────────────────────────┐
                              │                         │                         │
                              ▼                         ▼                         ▼
                     ┌────────────────┐       ┌────────────────┐       ┌────────────────┐
                     │    Worker 1    │       │    Worker 2    │       │    Worker N    │
                     │                │       │                │       │                │
                     │ Execute Task   │       │ Execute Task   │       │ Execute Task   │
                     └───────┬────────┘       └───────┬────────┘       └───────┬────────┘
                             │                        │                        │
                             └────────────────────────┼────────────────────────┘
                                                      │
                                                      ▼
                                           ┌──────────────────────┐
                                           │   Task State Manager │
                                           │                      │
                                           │ PENDING              │
                                           │ READY                │
                                           │ RUNNING              │
                                           │ SUCCESS              │
                                           │ FAILED               │
                                           │ RETRY                │
                                           └──────────┬───────────┘
                                                      │
                                                      ▼
                                           ┌──────────────────────┐
                                           │ Failure / Recovery   │
                                           │                      │
                                           │ Retries              │
                                           │ Timeouts             │
                                           │ Worker Failure       │
                                           │ Task Recovery        │
                                           └──────────────────────┘
```

----------

# Final Execution Flow

```text
Workflow Submitted
        │
        ▼
   Validate DAG
        │
        ▼
  Detect Cycles
        │
        ▼
  Initialize Tasks
        │
        ▼
 Find Ready Tasks
        │
        ▼
Apply Scheduling Policy
        │
        ▼
   Priority Queue
        │
        ▼
 Assign Tasks to Workers
        │
        ▼
 Concurrent Execution
        │
        ▼
  Track Task State
        │
        ├──────────────────────┐
        │                      │
        ▼                      ▼
     SUCCESS                FAILED
        │                      │
        │                   Retry?
        │                  /       \
        │                YES        NO
        │                 │          │
        │                 ▼          ▼
        │               Retry      FAILED
        │
        ▼
Find Newly Ready Tasks
        │
        ▼
 Schedule Again
        │
        ▼
Workflow Complete
```

----------

# Distributed Worker Model

The final scheduler will support multiple workers executing independent  
tasks concurrently.

```text
                         Scheduler
                             │
                ┌────────────┼────────────┐
                │            │            │
                ▼            ▼            ▼
             Worker 1     Worker 2     Worker 3
                │            │            │
                ▼            ▼            ▼
             Task A       Task B       Task C
                │            │            │
                └────────────┼────────────┘
                             │
                             ▼
                       Task Completion
```

Workers will eventually need to support:

```text
Task Claim
    ↓
Task Execution
    ↓
Heartbeat
    ↓
Task Completion
```

Worker failure:

```text
Worker crashes
      ↓
Task remains unfinished
      ↓
Scheduler detects failure
      ↓
Task becomes eligible again
      ↓
Another worker claims task
      ↓
Task executes
```

----------

# Core Concepts

## 1. Tasks

A task represents a unit of work.

```go
type Task struct {
    Name     string
    Status   Status
    Priority int
}
```

Current lifecycle:

```text
PENDING
   ↓
READY
   ↓
RUNNING
   ↓
SUCCESS
```

A task may also become:

```text
FAILED
```

The task itself stores execution state and metadata.

Dependencies are intentionally kept in the graph rather than inside the  
task.

----------

# 2. Dependency Graph

The graph stores relationships between tasks.

Example:

```text
        A
       / \
      B   C
       \ /
        D
```

Represented as:

```text
A → B
A → C
B → D
C → D
```

The graph can answer:

```text
Who are the parents of a task?

Who are the children of a task?
```

For example:

```go
g.Parents("D")
```

returns:

```text
[B C]
```

while:

```go
g.Children("A")
```

returns:

```text
[B C]
```

The graph is responsible for dependency relationships.

----------

# 3. Ready Task Detection

The scheduler determines which tasks can currently execute.

A task is ready when:

1.  It is still `PENDING`.
    
2.  It has no dependencies, or
    
3.  All of its dependencies have completed successfully.
    

Example:

```text
        A
       / \
      B   C
       \ /
        D
```

Initially:

```text
A → READY
B → WAITING
C → WAITING
D → WAITING
```

After `A` succeeds:

```text
B → READY
C → READY
D → WAITING
```

After both `B` and `C` succeed:

```text
D → READY
```

This is the core mechanism that allows the scheduler to move through a DAG.

----------

# 4. Priority Scheduling

Multiple tasks can become ready at the same time.

Example:

```text
B → priority 5
C → priority 10
```

Both are ready.

The scheduler uses the priority queue:

```text
Ready Tasks
     │
     ▼
Priority Queue
     │
     ▼
C (10)
     │
     ▼
B (5)
```

Therefore:

```text
C executes before B
```

Priority allows the scheduler to make an explicit scheduling decision instead  
of relying on task insertion order.

----------

# 5. Priority Queue

The current priority queue is implemented manually.

Current complexity:

```text
Push → O(1)
Pop  → O(n)
```

The queue can later be replaced with a heap-based implementation:

```text
Push → O(log n)
Pop  → O(log n)
```

This provides a natural opportunity to understand priority queues and heaps  
instead of hiding the implementation behind a library.

----------

# 6. Concurrent Execution

The next major stage is concurrent execution.

If:

```text
        A
       / \
      B   C
```

and `A` finishes, both `B` and `C` can execute concurrently.

Sequential execution:

```text
B █████████
           C █████████
```

Concurrent execution:

```text
B █████████
C █████████
```

The scheduler will gradually introduce:

-   Goroutines
    
-   `sync.WaitGroup`
    
-   Channels
    
-   Worker pools
    
-   Concurrency limits
    
-   Synchronization
    
-   Race-condition handling
    

The goal is to understand why independent DAG branches can safely execute at  
the same time.

----------

# 7. Failure Handling

Eventually tasks will be able to fail.

Example:

```text
Task
 ↓
FAILED
 ↓
Retry?
 ├── YES → RUNNING
 │           ↓
 │        SUCCESS
 │
 └── NO → FAILED
```

Planned capabilities:

-   Retry attempts
    
-   Maximum retry count
    
-   Retry delays
    
-   Task timeout
    
-   Cancellation
    
-   Failure propagation
    
-   Recovery
    

----------

# 8. Distributed Execution

After local concurrency works, the scheduler will evolve into a distributed  
execution model.

```text
                   Scheduler
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
       Worker 1     Worker 2     Worker 3
```

This introduces new problems:

-   Task claiming
    
-   Worker identity
    
-   Worker heartbeats
    
-   Worker failure
    
-   Duplicate execution
    
-   Task ownership
    
-   Recovery
    
-   Idempotency
    
-   Coordination
    

The purpose is to understand the additional problems introduced when execution  
moves from one process to multiple workers.

----------

# Development Roadmap

## Phase 1 — Core DAG Scheduler

-   Task structure
    
-   Task states
    
-   Dependency graph
    
-   Add tasks
    
-   Add dependencies
    
-   Parent lookup
    
-   Child lookup
    
-   Ready task detection
    
-   Basic sequential scheduler
    

----------

## Phase 2 — Priority Scheduling

-   Priority field
    
-   Priority queue
    
-   Queue push
    
-   Queue pop
    
-   Highest-priority task selection
    
-   Integrate priority queue with scheduler
    

----------

## Phase 3 — Concurrent Scheduler

-   Understand goroutines
    
-   `sync.WaitGroup`
    
-   Concurrent execution of independent tasks
    
-   Synchronize task completion
    
-   Worker pool
    
-   Concurrency limits
    
-   Channels
    
-   Race-condition testing
    

----------

## Phase 4 — Reliable Scheduler

-   Task failures
    
-   Retries
    
-   Retry policies
    
-   Timeouts
    
-   Cancellation
    
-   Failure propagation
    
-   Recovery
    

----------

## Phase 5 — Distributed Scheduler

-   Multiple workers
    
-   Task claiming
    
-   Worker identity
    
-   Worker heartbeats
    
-   Worker failure detection
    
-   Task reassignment
    
-   Duplicate execution handling
    
-   Idempotent execution
    
-   Distributed coordination
    

----------

## Phase 6 — Advanced Scheduling

Eventually the scheduling policy can consider more than priority:

```text
Priority
   +
SLA
   +
Task Age
   +
Estimated Execution Time
   +
Business Importance
   +
Worker Availability
```

This allows experimentation with different scheduling strategies.

----------

# Project Structure

```text
mini-scheduler/
│
├── task/
│   └── task.go
│
├── graph/
│   └── graph.go
│
├── queue/
│   └── priority_queue.go
│
├── scheduler/
│   └── scheduler.go
│
├── main.go
│
├── README.md
│
└── go.mod
```

As the project evolves, additional packages will be introduced for workers,  
execution, reliability, and distributed coordination.

----------

# Example Workflow

Consider:

```text
                 Orders
                /      \
               /        \
          Payments     Inventory
               \        /
                \      /
                 Revenue
```

The scheduler initially executes:

```text
Orders
```

After `Orders` succeeds:

```text
Payments
Inventory
```

become ready.

If:

```text
Payments   → priority 10
Inventory  → priority 5
```

priority scheduling selects:

```text
Payments
```

first.

Once both upstream tasks complete:

```text
Revenue
```

becomes ready.

Eventually, with concurrency:

```text
             Orders
                │
        ┌───────┴────────┐
        ▼                ▼
    Payments          Inventory
        │                │
        └───────┬────────┘
                ▼
             Revenue
```

`Payments` and `Inventory` can execute concurrently.

----------

# Purpose of building this project

The project is designed to build a practical understanding of:

### DAGs

-   Dependency representation
    
-   Parent/child relationships
    
-   Ready-node detection
    
-   Topological execution
    
-   Cycle detection
    

### Scheduling

-   Priority scheduling
    
-   Priority queues
    
-   Scheduling policies
    
-   Fairness
    
-   Starvation
    
-   Concurrency limits
    

### Concurrency

-   Goroutines
    
-   Channels
    
-   Synchronization
    
-   Worker pools
    
-   Race conditions
    
-   Task coordination
    

### Reliability

-   Retries
    
-   Timeouts
    
-   Failure recovery
    
-   Idempotency
    
-   Cancellation
    

### Distributed Systems

-   Worker coordination
    
-   Task ownership
    
-   Heartbeats
    
-   Failure detection
    
-   Reassignment
    
-   Duplicate execution
    
-   Distributed task state
    

----------

# Design Principles

## Build From First Principles

Core mechanisms are implemented manually wherever practical.

Instead of immediately using a workflow framework, the project first implements  
the underlying concepts directly.

----------

## Separate Responsibilities

```text
Task
    ↓
State + Metadata

Graph
    ↓
Dependencies

Priority Queue
    ↓
Scheduling Order

Scheduler
    ↓
Orchestration

Worker
    ↓
Execution
```

Each component should have a clear responsibility.

----------

## Increase Complexity Gradually

The system evolves through deliberate stages:

```text
Sequential Scheduler
        ↓
Priority Scheduler
        ↓
Concurrent Scheduler
        ↓
Reliable Scheduler
        ↓
Distributed Scheduler
```

Each stage introduces a real engineering problem.

----------

# Testing

Run all tests:

```bash
go test ./...
```

Run the scheduler:

```bash
go run main.go
```

Run with the race detector once concurrency is introduced:

```bash
go run -race main.go
```




```


