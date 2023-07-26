import MacroManager from "./components/MacroManager";
import TokenManager from "./components/TokenManager";

export default function App() {
    return (
        <main className="container">
            <hgroup>
                <h1>One Time Service</h1>
                <h2>Internal Dashboard</h2>
            </hgroup>
            <MacroManager />
            <TokenManager />
        </main>
    )
}
