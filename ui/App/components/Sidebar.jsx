import React from 'react';
import {Link, IndexLink} from 'react-router';

class Sidebar extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        if (this.props.serverRunning === "running") {
            var serverStatus = 
                <IndexLink to="/"><i className="fa fa-circle text-success" />Servers Online</IndexLink>
        } else {
            var serverStatus = 
                <IndexLink to="/"><i className="fa fa-circle text-danger" />Servers Offline</IndexLink>
        }

        return(
            <aside className="main-sidebar">

                <section className="sidebar" style={{height: "100%"}}>

                <div className="user-panel">
                    <div className="pull-left image">
                    <img src="./dist/dist/img/factorio.jpg" className="img-circle" alt="User Image" />
                    </div>
                    <div className="pull-left info">
                    <p>Factorio Server Manager</p>

                    {serverStatus}

                    </div>
                </div>

                <form action="#" method="get" className="sidebar-form">
                    <div className="input-group">
                    <input type="text" name="q" className="form-control" placeholder="Search..." />
                        <span className="input-group-btn">
                            <button type="submit" name="search" id="search-btn" className="btn btn-flat"><i className="fa fa-search" />
                            </button>
                        </span>
                    </div>
                </form>

                <ul className="sidebar-menu">
                    <li className="header">MENU</li>
                    <li><IndexLink to="/" activeClassName="active"><i className="fa fa-tachometer" /><span>Server Control</span></IndexLink></li>
                    <li><Link to="/configs" activeClassName="active"><i className="fa fa-cogs" /><span>Server Configs</span></Link></li>
                    <li><Link to="/logs" activeClassName="active"><i className="fa fa-file-text-o" /> <span>Logs</span></Link></li>
                    <li><Link to="/saves" activeClassName="active"><i className="fa fa-floppy-o" /> <span>Saves</span></Link></li>
                    <li><Link to="/config" activeClassName="active"><i className="fa fa-cogs" /> <span>Game Configuration</span></Link></li>
                    <li><Link to="/settings" activeClassName="active"><i className="fa fa-cog" /> <span>Settings</span></Link></li>
                </ul>
                </section>
                <div style={{height: "100%"}} />
            </aside>
        )
    }
}

Sidebar.propTypes = {
    
}

export default Sidebar
