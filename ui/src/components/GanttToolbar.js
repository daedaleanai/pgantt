
import React, { Component } from 'react';
import { PageHeader, Radio } from 'antd';

class GanttToolbar extends Component {
  state = {
    currentZoom: 'Days',
  };

  handleZoomChange = (e) => {
    this.setState({
      currentZoom: e.target.value,
    });

    if (this.props.onZoomChange) {
      this.props.onZoomChange(e.target.value);
    }
  }

  render() {
    const options = [
      { label: 'Days', value: 'Days' },
      { label: 'Months', value: 'Months' },
    ];

    return (
      <div className="zoom-bar">
        <PageHeader
          ghost={false}
          title="Platforms"
          extra={[
            <Radio.Group
              options={options}
              onChange={this.handleZoomChange}
              value={this.state.currentZoom}
              optionType="button"
            />,
          ]}
        >
        </PageHeader>
      </div>
    );
  }
}

export default GanttToolbar;
