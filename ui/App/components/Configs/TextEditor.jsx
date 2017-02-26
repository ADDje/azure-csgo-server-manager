import React from 'react'

class TextEditor extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            isLoading: false,
            content: "",
            error: null
        }

        this.loadContent = this.loadContent.bind(this)
        this.changeText = this.changeText.bind(this)
        this.save = this.save.bind(this)
    }

    componentWillMount() {
        this.loadContent(this.props.name, this.props.type, this.props.template)
    }

    componentWillReceiveProps(nextProps) {
        if (nextProps.name !== this.props.name ||
            nextProps.type !== this.props.type) {
            // Populate the parameters
            this.loadContent(nextProps.name, nextProps.type, nextProps.template)
        }
    }

    changeText(e) {
        this.setState({content: e.target.value})
    }

    loadContent(name, type, data) {
        this.setState({isLoading: true})

        if (name === undefined) {
            name = this.props.name
        }
        if (type === undefined) {
            type = this.props.type
        }

        var url
        switch (type) {
            case 'config':
                url = "/api/configs/gettext/" + name
                break
            case 'template':
                this.setState({content: JSON.stringify(data, null, 4), isLoading: false})
                return
            default:
                console.log("Invalid TextEditor Content: " + type)
        }

        $.ajax({
            type: "GET",
            url: url,
            dataType: "json",
            success: (resp) => {
                this.setState({
                    content: resp.data,
                    isLoading: false
                })
            }
        })
    }

    save() {
        var url
        switch (this.props.type) {
            case 'config':
                url = "/api/configs/gettext/" + this.props.name
                break
            case 'template':
                url = "/api/templates/" + this.props.name + "/update"
                break
            default:
                console.log("Invalid TextEditor Content: " + this.props.type)
                return
        }

        this.setState({isLoading: true})

        $.ajax({
            type: "POST",
            url: url,
            dataType: "json",
            data: this.state.content,
            success: (resp) => {
                if (typeof(resp.success) === "undefined" || resp.success === false) {
                    this.setState({isLoading: false, error: resp.data})
                } else {
                    this.setState({isLoading: false, error: null})
                    
                    if (this.props.reloadSelected !== null) {
                        this.props.reloadSelected()
                    }
                }
            }
        })
    }

    render() {
        
        if (this.props.name === null || this.props.name === undefined) {
            return null
        }

        var overlay = null
        if (this.state.isLoading) {
            overlay = (<div className="overlay">
                <i className="fa fa-refresh fa-spin" />
            </div>)
        }

        var errorMessage = null
        if (this.state.error) {
            errorMessage = (<div className="callout callout-danger">
                <h4><i className="icon fa fa-ban" />Error!</h4>
                {this.state.error}
              </div>)
        }

        return(
            <div className="box">
                {errorMessage}
                <div className="text-editor">
                    <textarea className="full-text" value={this.state.content} onChange={this.changeText} />
                    <button onClick={this.save} className="btn btn-primary">Update</button>
                </div>
                {overlay}
            </div>
        )
    }
}

TextEditor.propTypes = {
    reloadSelected: React.PropTypes.func.isRequired,
    type: React.PropTypes.string.isRequired,
}

export default TextEditor
