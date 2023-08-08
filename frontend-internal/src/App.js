import React, { useState, useEffect } from "react";
import MacroManager from "./components/MacroManager";
import TokenManager from "./components/TokenManager";

export default function App() {
    const [ macros, setMacros ] = useState(null)
    const [ tokens, setTokens ] = useState(null)
    
    useEffect(() => {
        fetch("api/internal/tokens")
            .then(r => r.json())
            .then(r => {
                setTokens(r.response)
            })
            .catch(console.error)

        fetch("api/internal/macros")
            .then(r => r.json())
            .then(r => {
                setMacros(r.response)
            })
            .catch(console.error)
    }, [])

    return (
        <main className="container">
            <br />
            <hgroup>
                <h1>One Time Service</h1>
                <h2>Internal Dashboard</h2>
            </hgroup>
            <MacroManager 
                macros={macros}
            />
            <TokenManager 
                macros={macros}
                tokens={tokens}
            />
        </main>
    )
}
