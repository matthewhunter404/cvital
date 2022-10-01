import React, { useState } from 'react'
import './App.css'
import axios from 'axios'

function LoginForm() {
    const [errorMessages, setErrorMessages] = useState({})
    const [isSubmitted, setIsSubmitted] = useState(false)

    const errors = {
        invalid: 'invalid username or password',
    }

    // Generate JSX code for error message
    const renderErrorMessage = (name) =>
        name === errorMessages.name && (
            <div className="error">{errorMessages.message}</div>
        )

    const handleSubmit = (event) => {
        // Prevent page reload
        event.preventDefault()
        var { user, pass } = document.forms[0]
        console.log('userValue', user.value)
        console.log('passValue', pass.value)
        const config = {
            method: 'get',
            url: 'localhost/user/login',
            port: 3000,
        }

        axios
            .post(
                'http://localhost:3000/user/login',
                {
                    email: user.value,
                    password: pass.value,
                },
                {
                    headers: {
                        'Content-Type': 'application/json',
                    },
                }
            )
            .then((response) => {
                console.log(response.data)
                setIsSubmitted(true)
            })
            .catch(function (error) {
                console.log('Oh no an error', error)
                setErrorMessages({ name: 'invalid', message: errors.invalid })
            })
    }

    // JSX code for login form
    const renderForm = (
        <div className="form">
            <form onSubmit={handleSubmit}>
                <div className="input-container">
                    <label>Email </label>
                    <input type="text" name="user" required />
                </div>
                <div className="input-container">
                    <label>Password </label>
                    <input type="password" name="pass" required />
                    {renderErrorMessage('invalid')}
                </div>
                <div className="button-container">
                    <input type="submit" />
                </div>
            </form>
        </div>
    )

    return (
        <div className="login-form">
            <div className="title">Sign In</div>
            {isSubmitted ? (
                <div>User is successfully logged in</div>
            ) : (
                renderForm
            )}
        </div>
    )
}

export default LoginForm
