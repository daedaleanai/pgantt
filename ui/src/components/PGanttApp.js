
import React, { Component } from 'react';

import PGanttNav from './PGanttNav';
import Gantt from './Gantt';
import GanttToolbar from './GanttToolbar';

class PGanttApp extends Component {
  state = {
    currentZoom: 'Days'
  };

  handleZoomChange = (zoom) => {
    this.setState({
      currentZoom: zoom
    });
  }

  render() {
    const { currentZoom } = this.state;
    return (
      <div>
        <PGanttNav />
        <GanttToolbar
          zoom={currentZoom}
          onZoomChange={this.handleZoomChange}
        />
        <Gantt zoom={currentZoom} />
      </div>
    );
  }
}

export default PGanttApp;
