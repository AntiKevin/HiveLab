import {
    Quit,
    WindowMinimise,
    WindowToggleMaximise,
} from "../../wailsjs/runtime/runtime";
import "./TopBar.css";

function TopBar() {
    return (
        <header className="top-bar">
            <div className="window-drag-region">
                <div className="app-title">
                    <span>SmokeLab</span>
                </div>
            </div>

            <div className="window-actions" aria-label="Window controls">
                <button type="button" className="window-action" onClick={WindowMinimise} aria-label="Minimize">
                    <span aria-hidden="true"></span>
                </button>
                <button type="button" className="window-action" onClick={WindowToggleMaximise} aria-label="Maximize">
                    <span aria-hidden="true"></span>
                </button>
                <button type="button" className="window-action close" onClick={Quit} aria-label="Close">
                    <span aria-hidden="true"></span>
                </button>
            </div>
        </header>
    );
}

export default TopBar;
