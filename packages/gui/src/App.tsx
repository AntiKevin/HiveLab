import './App.css';
import TopBar from "./components/TopBar";

function App() {
    return (
        <div className="app-shell">
            <TopBar />

            <main className="workspace">
                <section className="empty-state" aria-label="Workspace">
                    <div className="empty-state-line"></div>
                    <h1>HiveLab</h1>
                </section>
            </main>
        </div>
    );
}

export default App;
