import React from 'react';
import update from 'immutability-helper';

class ScheduleEditor extends React.Component {

    constructor(props) {
        super(props)

        this.state = {
            name: "",
            action: "",
            parameters: [],
            enabled: true,

            nameError: false,
            actionError: false
        }

        this.changeField = this.changeField.bind(this)
        this.save = this.save.bind(this)

        this.changeKey = this.changeKey.bind(this)
        this.changeValue = this.changeValue.bind(this)
    }

    componentWillReceiveProps(nextProps) {
        if (nextProps.selectedActionName !== this.props.selectedActionName) {
            if (nextProps.selectedActionName !== "" &&
                nextProps.actions[nextProps.selectedActionName] !== undefined) {
                
                var action = nextProps.actions[nextProps.selectedActionName]
                this.setState({
                    name: action.name,
                    action: action.action,
                    parameters: action.parameters,
                    enabled: action.enabled,
                    
                    nameError: false,
                    actionError: false
                })
            } else {
                this.setState({
                    name: "",
                    action: "",
                    parameters: [],
                    enabled: true,

                    nameError: false,
                    actionError: false
                })
            }
        }
    }

    changeKey(key, e) {
        if (typeof(e) === "undefined") {
            e = key

            this.setState({
                parameters: update(this.state.parameters, {$push: [{key: e.target.value, value: ""}]})
            })
        } else {

            var params = this.state.parameters.slice()
            if (e.target.value === "" && this.state.parameters[key].value === "") {
                delete params[key]
            } else {
                params[key] = {key: e.target.value, value: this.state.parameters[key].value}
            }
            this.setState({
                parameters: params
            })
        }
    }

    changeValue(key, e) {
        if (typeof(e) === "undefined") {
            e = key

            this.setState({
                parameters: update(this.state.parameters, {$push: [{value: e.target.value, key: ""}]})
            })
        } else {

            var params = this.state.parameters.slice()
            if (e.target.value === "" && this.state.parameters[key].key === "") {
                delete params[key]
            } else {
                params[key] = {value: e.target.value, key: this.state.parameters[key].key}
            }
            this.setState({
                parameters: params
            })
        }
    }

    changeField(key, event) {
        
        var val = (key === "enabled") ? event.target.checked : event.target.value

        var param = {}
        param[key] = {$set: val}
        this.setState(update(this.state, param))
    }

    save(e) {
        e.preventDefault()

        var r = /[^a-zA-Z0-9-_]/
        var nameError = false
        var actionError = false
        if (this.state.name.length < 1 || this.state.name.match(r)) {
            nameError = true
        }
        if (this.state.action === "") {
            actionError = true
        }
        this.setState({
            nameError: nameError,
            actionError: actionError
        })
        if (nameError || actionError) {
            return
        }

        var data = {
            name: this.state.name,
            action: this.state.action,
            parameters: this.state.parameters,
            enabled: this.state.enabled
        }

        var name = (this.props.selectedActionName === "") ? this.state.name : this.props.selectedActionName
        $.post({
            url: "/api/schedule/" + name,
            dataType: "json",
            data: JSON.stringify(data, null, 4),
            success: (resp) => {
                if (typeof(resp.success) === "undefined" || resp.success === false) {
                    this.setState({isLoading: false, error: resp.data});
                } else {
                    this.setState({isLoading: false, error: null});
                    this.props.reloadSelected()
                }
            }
        })
    }

    render() {
        var method = (this.props.selectedActionName === "") ? "Create" : "Edit"
        var title = method + " Schedule Action"

        var parameters = []
        for (var p in this.state.parameters) {
            var name = "parameter-" + p
            var keyName = name + "-key"
            var valueName = name + "-value"
            parameters.push(<div key={name} id={name} className="input-group">
                <input ref={keyName} name={keyName} id={keyName} type="text" onChange={this.changeKey.bind(this, p)} value={this.state.parameters[p].key} />
                <input ref={valueName} name={valueName} id={valueName} type="text" onChange={this.changeValue.bind(this, p)} value={this.state.parameters[p].value} />
            </div>)
        }
        var name = "parameter-" + this.state.parameters.length
        var keyName = name + "-key"
        var valueName = name + "-value"
        parameters.push(<div key={name} id={name} className="input-group">
            <input ref={keyName} name={keyName} id={keyName} type="text" onChange={this.changeKey} value="" />
            <input ref={valueName} name={valueName} id={valueName} type="text" onChange={this.changeValue} value="" />
        </div>)

        var actionHelp, nameHelp = null
        var actionClass, nameClass = "input-group"
        
        if (this.state.actionError) {
            actionHelp = <span className="help-block">Action is required</span>
            actionClass += " has-error"
        }
        if (this.state.nameError) {
            nameHelp = <span className="help-block">Invalid name. Can only contain characters allowed in urls</span>
            nameClass += " has-error"
        }

        return (
            <div className="template-parameters box box-primary">
                <div className="box-header with-border">
                    <h3 className="box-title">{title}</h3>
                </div>
                
                <form role="form">
                    <div className="box-body">
                        <label htmlFor="name">Name (slug)</label>
                        <div id="name" className={nameClass}>
                            <input ref="name" name="name" id="name" type="text" className="form-control" onChange={this.changeField.bind(this, "name")} value={this.state.name} />
                        </div>
                        {nameHelp}

                        <label htmlFor="action">Action</label>
                        <div id="action" className={actionClass}>
                            <select ref="action" name="action" id="action" className="form-control" onChange={this.changeField.bind(this, "action")} value={this.state.action} >
                                <option />
                                <option>Deploy</option>
                                <option>Delete</option>
                                <option>Start</option>
                                <option>Stop</option>
                            </select>
                        </div>
                        {actionHelp}


                        <label htmlFor="params">Parameters</label>
                        {parameters}

                        <label htmlFor="name">Enabled</label>
                        <div id="enabled" className="input-group">
                            <input ref="enabled" name="enabled" id="enabled" type="checkbox" onChange={this.changeField.bind(this, "enabled")} checked={this.state.enabled} />
                        </div>
                            
                        <div style={{marginTop: 10 + 'px'}}>
                            <button onClick={this.save} type="submit" className="btn btn-primary">{method}</button>
                        </div>
                    </div>
                </form>
            </div>
        );
    }

}

ScheduleEditor.propTypes = {
    reloadSelected: React.PropTypes.func.isRequired,
}

export default ScheduleEditor