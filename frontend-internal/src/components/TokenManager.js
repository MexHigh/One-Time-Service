import React from "react"
import TokenAdder from "./TokenAdder"
import TokenEntry from "./TokenEntry"

export default function TokenManager({ serviceCalls, tokens }) {
    return (
        <article>
            <header>
                Manage tokens
            </header>
            { !serviceCalls ? (
                <p>Please add a service call first!</p>
            ) : (
                <div>
                    <TokenAdder serviceCalls={serviceCalls} />
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
            )}
        </article>
    )
}