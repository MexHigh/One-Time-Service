import React from "react"
import TokenAdder from "./TokenAdder"
import TokenEntry from "./TokenEntry"

export default function TokenManager({ macros, tokens }) {
    return (
        <>
            <article>
                <header>
                    Manage <span
                        data-tooltip="One time code that will call a Macro once executed"
                        data-placement="right"
                    >Tokens</span>
                </header>
                <div>
                    <TokenAdder macros={macros} />
                    <div>
                        { tokens && Object.entries(tokens).map(([token, details]) => (
                            <TokenEntry
                                key={token}
                                token={token}
                                details={details}
                            />
                        ))}
                    </div>
                </div>
                {/*<footer>
                    <button className="outline">
                        Delete expired Tokens
                    </button>
                    <button className="outline">
                        Delete all Tokens
                    </button>
                </footer>*/}
            </article>
        </>
    )
}