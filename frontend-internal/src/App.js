import React, { useState, useEffect } from "react";
import ServiceCallManager from "./components/ServiceCallManager";
import TokenManager from "./components/TokenManager";

export default function App() {
    const [ serviceCalls, setServiceCalls ] = useState(null)
    const [ tokens, setTokens ] = useState(null)
    
    useEffect(() => {
        fetch("api/internal/tokens/details")
            .then(r => r.json())
            .then(r => {
                setTokens(r.response)
            })
            .catch(console.error)

        fetch("api/internal/service-calls")
            .then(r => r.json())
            .then(r => {
                setServiceCalls(r.response)
            })
            .catch(console.error)
    }, [])

    return (
        <main className="container">
            <br />
            <hgroup>
                <h1>One Time Service</h1>
                <h2>Internal dashboard</h2>
            </hgroup>
            <ServiceCallManager 
                serviceCalls={serviceCalls}
            />
            <TokenManager 
                serviceCalls={serviceCalls}
                tokens={tokens}
            />
        </main>
    )
}
