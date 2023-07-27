import React, { useEffect, useState } from "react"
import MacroAdder from "./MacroAdder"
import MacroModal from "./MacroModal"

export default function MacroManager() {
    const [ macros, setMacros ] = useState(null)
    const [ modal, setModal ] = useState(null)
    
    useEffect(() => {
        fetch("/api/internal/macros")
            .then(r => r.json())
            .then(r => {
                setMacros(r.response)
            })
            .catch(console.error)
    }, [])

    return (
        <>
            <article>
                <header>
                    Manage <span 
                        data-tooltip="Macros map a name to a HA service call"
                        data-placement="right"
                    >Macros</span>
                </header>
                <div>
                    <MacroAdder />
                    <h5>Macro List</h5>
                    <ul>
                        { macros && macros.map(macro => (
                            <li key={ macro }>
                                <a 
                                    onClick={e => {
                                        e.preventDefault()
                                        setModal(macro)
                                    }}
                                >{ macro }</a>
                            </li>
                        ))}
                    </ul>
                </div>
                <footer>
                    
                </footer>
            </article>
            <MacroModal
                open={modal ? true : false}
                closeCallback={() => setModal(null)}
                macroName={modal}
            />
        </>
    )
}