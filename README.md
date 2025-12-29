# 🔐 MatLock

![Go](https://img.shields.io/badge/Go-1.22+-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Educational-orange)

`MatLock` is a **simple mutex implementation written in Go**, built to help understand **how mutual exclusion works internally**, using channels and basic state instead of `sync.Mutex`.

This project is **educational** and demonstrates:

- What a mutex is  
- Why we need it  
- How it can be implemented using channels  

---

## 📚 Table of Contents

- [What is a Mutex?](#-what-is-a-mutex)
- [Using sync.Mutex](#-example-using-syncmutex)
- [How a Mutex Works Internally](#-how-does-a-mutex-work-internally)
- [MatLock Implementation](#-matlock-a-simple-mutex-implementation)
- [How It Works](#-how-it-works)
- [Why This Project?](#-why-this-project)
- [Disclaimer](#-disclaimer)

---

## 🧠 What is a Mutex?

A **Mutex (Mutual Exclusion)** is a locking mechanism used to ensure that **only one goroutine can access a critical section of code at a time**.

This prevents **race conditions**, especially when multiple goroutines read/write shared data.

In Go, the standard mutex lives in the `sync` package and provides two main methods:

- `Lock()`
- `Unlock()`

---

## ✅ Example: Using `sync.Mutex`

Below is a classic example of using `sync.Mutex` to safely access a shared map.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	// Only one goroutine can access c.v at a time
	c.v[key]++
	c.mu.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v[key]
}

func main() {
	c := SafeCounter{v: make(map[string]int)}

	for i := 0; i < 1000; i++ {
		go c.Inc("somekey")
	}

	time.Sleep(time.Second)
	fmt.Println(c.Value("somekey"))
}
```

---

## 🔍 How Does a Mutex Work Internally?

Conceptually, a mutex works like this:

- If the lock is **free**, a goroutine can enter the critical section  
- If the lock is **held**, other goroutines must **wait**  
- When the lock is released, **one waiting goroutine is notified**

You can think of a mutex as a coordination mechanism similar to a channel:

- When access is unavailable → goroutines block  
- When access becomes available → exactly one goroutine continues

This guarantees **mutual exclusion** over shared resources.

---

## 🧪 MatLock: A Simple Mutex Implementation

`MatLock` is a **minimal mutex implementation** created for learning purposes.

It is built using:

- A boolean flag (`locked`) to represent lock state  
- A buffered channel (`wait`) to block and wake goroutines  

⚠️ **Important:**  
This implementation is **not intended for production use**.  
Always use `sync.Mutex` in real-world applications.

---

### Implementation

```go
package matlock

type MatLock struct {
	locked bool
	wait   chan struct{}
}

func New() *MatLock {
	return &MatLock{
		locked: false,
		wait:   make(chan struct{}, 1),
	}
}

func (m *MatLock) Lock() {
	for {
		if !m.locked {
			m.locked = true
			return
		}
		// Wait until the lock is released
		<-m.wait
	}
}

func (m *MatLock) Unlock() {
	if !m.locked {
		panic("unlock of unlocked MatLock")
	}

	m.locked = false

	// Wake up exactly one waiting goroutine (if any)
	select {
	case m.wait <- struct{}{}:
	default:
	}
}
```

---

## 🛠 How It Works

1. The `locked` flag represents whether the critical section is currently occupied  
2. When `Lock()` is called:
   - If `locked == false`, the goroutine acquires the lock
   - If `locked == true`, the goroutine blocks on the `wait` channel
3. When `Unlock()` is called:
   - The lock is released by setting `locked = false`
   - One waiting goroutine (if any) is notified via the channel
4. This ensures that **only one goroutine** can enter the critical section at any time

---

## 🎯 Why This Project?

This project was created to:

- Understand how mutexes work internally  
- Learn Go concurrency primitives (goroutines & channels)  
- Explore synchronization without relying on `sync.Mutex`  
- Build intuition about race conditions and mutual exclusion  

It is intentionally **simple and minimal** to focus on the core idea.

---

## 📌 Disclaimer

This implementation:

- Is **not thread-safe** at the CPU / Go memory model level  
- Does **not guarantee fairness** between goroutines  
- Does **not handle starvation or priority**  
- Is **not a replacement for `sync.Mutex`**  

⚠️ **Do not use this in production code.**  
It exists **only for learning and experimentation**.
