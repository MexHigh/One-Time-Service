import React, { useEffect, useState } from "react"
import MacroAdder from "./MacroAdder"

export default function MacroManager() {
    const [ macros, setMacros ] = useState(null)
    
    useEffect(() => {
        fetch("/api/internal/macros")
            .then(r => r.json())
            .then(r => {
                setMacros(r.response)
            })
            .catch(console.error)
    }, [])

    return (
        <article>
            <header>
                Manage Macros
            </header>
            <div>
                <MacroAdder />
                <h5>Macro List</h5>
                <ul>
                    { macros && macros.map(macro => (
                        <li key={ macro }>
                            <a>{ macro }</a>
                        </li>
                    ))}
                </ul>
            </div>
            <footer>
                
            </footer>
        </article>
    )
}