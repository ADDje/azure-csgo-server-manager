import React from 'react';
import {IndexLink} from 'react-router';
import UserTable from './Users/UserTable.jsx';
import AddUser from './Users/AddUser.jsx';

class ServerLog extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            expanded: false
        }

        this.expand = this.expand.bind(this)
    }

    componentDidMount() {

    }

    expand() {
        this.setState({
            expanded: !this.state.expanded
        })
    }

    render() {
        var mainClass = "collapse " + ((this.state.expanded) ? "in" : "");
        var summaryClass = "log-last-line collapse " + ((!this.state.expanded) ? "in" : "");
        var buttonClass = "fa btn-xs " + ((this.state.expanded) ? "fa-caret-square-o-down" : "fa-caret-square-o-up")

        return(
            <div className="log-footer">
                <a onClick={this.expand} href="#">
                    <div className={summaryClass}>
                        <i className="fa fa-tv" />
                        <p>Stuff happened</p>
                    </div>
                </a>
                <div className={mainClass} id="server-log">
                    <div className="server-log-title">
                        <a onClick={this.expand} href="#">
                            <i className="fa fa-tv" />
                            <h4>Server Manager Log</h4>
                        </a>
                    </div>
                    <div id="log">
                        <p>Something</p>
                        <p>Something else</p>
                    </div>
                </div>
            </div>
        )
    }
}

export default ServerLog
