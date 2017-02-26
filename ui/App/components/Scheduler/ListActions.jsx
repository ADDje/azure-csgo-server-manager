import React from 'react'

class ListActions extends React.Component {

    constructor(params) {
        super(params)

        this.canCreate = this.canCreate.bind(this)
        this.createClick = this.createClick.bind(this)

        this.deleteClick = this.deleteClick.bind(this)
        this.deleteAction = this.deleteAction.bind(this)

        this.clickRow = this.clickRow.bind(this)

        this.loadApiKey = this.loadApiKey.bind(this)
        this.getUri = this.getUri.bind(this)

        this.state = {
            apiKey: "",
            loadingKey: false
        }
    }

    componentWillMount() {
        if (this.state.apiKey === "") {
            this.loadApiKey()
        }
    }

    // TODO: This is inefficient. Settings should be shared, or only pull key
    loadApiKey() {
        this.setState({
            loadingKey: true
        })
        
        $.ajax({
            url: "/api/settings",
            dataType: "json",
            success: (resp) => {
                if (resp.success === true) {
                    this.setState({loadingKey: false, apiKey: resp.data.external_api_key})
                }
            },
            error: (xhr, status, err) => {
                console.log('/api/settings', status, err.toString())
            }
        })
    }

    canCreate() {
        return this.props.selectedActionName === ""
    }

    deleteClick(name, e) {
        e.preventDefault()
        swal({
            title: "Are you sure?",
            text: "You will not be able to recover " + name,
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Yes, delete it!",
            closeOnConfirm: false,
            showLoaderOnConfirm: true
        },
        function(){
            this.deleteAction(name)
        }.bind(this))
    }

    deleteAction(name) {
        $.post({
            url: "api/schedule/" + name + "/delete",
            dataType: "json",
            success: (resp) => {

                swal({
                    timer: 2000,
                    title: "Nice!",
                    text: name + " deleted!",
                    type: "success"
                })
                
                this.props.reloadActions()
                if (this.props.selectedActionName === name) {
                    this.props.selectAction("", null)
                }
            }
        })
    }

    createClick(e) {
        this.props.selectAction("", null)
    }

    clickRow(name, action, e) {
        e.preventDefault()
        if (e.target.nodeName === "INPUT") {
            return
        }
        this.props.selectAction(name, action)
    }

    clickSlug(e) {
        e.preventDefault()
        e.target.select()
    }

    getUri(action) {
        if (this.state.loadingKey) {
            return "Loading..."
        }
        return location.origin + "/external/schedule/" + action + "/exec?key=" + this.state.apiKey
    }

    render() {
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Manage Schedule Actions</h3>
                </div>
                
                <div className="box-body">
                    <div className="table-responsive">
                        <table className="table table-striped">
                            <thead>
                                <tr>
                                    <th width="20%">Name</th>
                                    <th width="30%" />
                                    <th width="25%">Action</th>
                                    <th width="25%">Enabled</th>
                                    <th width="5%" />
                                </tr>
                            </thead>
                            <tbody>
                            {Object.keys(this.props.actions).map ( (action, i) => {
                                var a = this.props.actions[action]
                                return(
                                    <tr key={action} onClick={this.clickRow.bind(this, action, a)}
                                        className={"row-link" + ((this.props.selectedActionName === action) ? " selected-row" : "")}
                                    >
                                        <td>
                                            {action}
                                        </td>
                                        <td>
                                            <input value={this.getUri(action)} readOnly onClick={this.clickSlug} />
                                        </td>
                                        <td>
                                            {a.action}
                                        </td>
                                        <td>
                                            {(a.enabled) ? "enabled" : "disabled"}
                                        </td>
                                        <td>
                                            <button className="btn btn-sm btn-danger" onClick={this.deleteClick.bind(this, action)}>
                                                <i className="fa fa-trash" />
                                            </button>
                                        </td>
                                    </tr>
                                )                                       
                            })}
                            </tbody>
                        </table>
                        <button className="btn btn-sm btn-success" disabled={this.canCreate()} onClick={this.createClick}>
                            <i className="fa fa-plus" />
                        </button>
                    </div>
                </div>
            </div>
        )
    }
}

ListActions.propTypes = {
    actions: React.PropTypes.object.isRequired,
    selectAction: React.PropTypes.func.isRequired,
    selectedActionName: React.PropTypes.string.isRequired,
}

export default ListActions
