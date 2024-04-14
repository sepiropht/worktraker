import React, { useState, useEffect } from 'react'
import './App.css'

function App() {
  const [tasks, setTasks] = useState([])
  const [description, setDescription] = useState('')
  const [day, setDay] = useState('')

  useEffect(() => {
    fetchTasks()
  }, [])

  const fetchTasks = () => {
    // Fetch tasks for the current day
    fetch(`/tasks`)
      .then((response) => response.json())
      .then((data) => setTasks(data))
      .catch((error) => console.error('Error fetching tasks:', error))
  }

  const addTask = () => {
    fetch('/addTask', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        description: description,
        day: day,
      }),
    })
      .then((response) => {
        if (response.ok) {
          fetchTasks()
          setDescription('')
        }
      })
      .catch((error) => console.error('Error adding task:', error))
  }

  const toggleTask = (id) => {
    fetch('/toggleTask', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: id,
        day: day,
      }),
    })
      .then((response) => {
        if (response.ok) {
          fetchTasks()
        }
      })
      .catch((error) => console.error('Error toggling task:', error))
  }

  const removeTask = (id) => {
    fetch('/removeTask', {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: id,
        day: day,
      }),
    })
      .then((response) => {
        if (response.ok) {
          fetchTasks()
        }
      })
      .catch((error) => console.error('Error removing task:', error))
  }

  return (
    <div className="App">
      <h1>Todo List</h1>
      <div>
        <input
          type="text"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Enter task description"
          required
        />
        <select value={day} onChange={(e) => setDay(e.target.value)} required>
          <option value="Monday">Monday</option>
          <option value="Tuesday">Tuesday</option>
          <option value="Wednesday">Wednesday</option>
          <option value="Thursday">Thursday</option>
          <option value="Friday">Friday</option>
          <option value="Saturday">Saturday</option>
          <option value="Sunday">Sunday</option>
        </select>
        <button onClick={addTask}>Add Task</button>
      </div>
      <h2>Tasks for Today ({day})</h2>
      <ul>
        {tasks.map((task) => (
          <li key={task.id}>
            <input
              type="checkbox"
              checked={task.done}
              onChange={() => toggleTask(task.id)}
            />
            <span className={task.done ? 'done' : ''}>{task.description}</span>
            <button onClick={() => removeTask(task.id)}>Remove</button>
          </li>
        ))}
      </ul>
    </div>
  )
}

export default App
