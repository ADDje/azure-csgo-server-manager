import React from 'react';
import {IndexLink} from 'react-router';
import UserTable from './Users/UserTable.jsx';
import AddUser from './Users/AddUser.jsx';
import update from 'immutability-helper';

class ServerLog extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            expanded: false,
            messages: []
        }

        this.expand = this.expand.bind(this)
        this.ws = null;
    }

    componentDidMount() {
        $.get({
            url: "/api/websocketinfo",
            success: (data) => {
                console.log(data)
                var server = (data.data.Address === "") ? window.location.hostname : data.data.server
                this.ws = new WebSocket("wss://" + server + ":" + data.data.Port)

                this.ws.onopen = function() {
                    console.log("Connected to log server")
                }

                this.ws.onmessage = function(msg) {
                    console.log(msg)
                    var splitData = msg.data.split("\n").filter(function(d) { return d !== "" })
                    console.log(splitData)
                    this.setState({messages: update(this.state.messages, {
                        $push: splitData
                    })})
                }.bind(this)
            }
        })
    }

    componentWillUnmount() {
        if (this.ws !== null) {
            this.ws.close();
        }
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

        var messages = this.state.messages
        var summary = (messages.length) ? (<p>{messages[messages.length-1]}</p>) : ""

        var i = 0;
        var fullMessages = messages.map(function(m) {
            return (<p key={i++}>{m}</p>)
        })

        return(
            <div className="log-footer">
                <a onClick={this.expand} href="#">
                    <div className={summaryClass}>
                        <i className="fa fa-tv" />
                        {summary}
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
                        {fullMessages}
                    </div>
                </div>
            </div>
        )
    }
}

export default ServerLog
