
import React, { Component } from 'react';

import PGanttNav from './PGanttNav';
import Gantt from './Gantt';

class PGanttApp extends Component {
  render() {
    return (
      <div>
        <PGanttNav />
        <Gantt />
      </div>
    );
  }
}

export default PGanttApp;
