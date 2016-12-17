import React from 'react';
import DynamicConfig from './DynamicConfig.jsx'
import TextEditor from './TextEditor.jsx'

class ConfigEditor extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        
        if (this.props.config === null) {
            return null
        }

        return(
            <div className="nav-tabs-custom">
                <ul className="nav nav-tabs">
                    <li><a href="#parameter_editor" data-toggle="tab">Parameter Editor</a></li>
                    <li><a href="#text_editor" data-toggle="tab">Text Editor</a></li>
                </ul>
                <div className="tab-content">
                    <div className="tab-pane" id="parameter_editor">
                        <DynamicConfig configName={this.props.configName} config={this.props.config} />
                    </div>
                    <div className="tab-pane" id="text_editor">
                        <TextEditor name={this.props.configName} type="config" reloadSelected={this.props.reloadSelected} />
                    </div>
                </div>
            </div>
        )
    }
}


export default ConfigEditor
