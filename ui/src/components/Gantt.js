
import React, { Component } from 'react';
import { gantt } from 'dhtmlx-gantt';

const data = {
    data: [
        { id: 1, text: 'Task #1', start_date: '15-04-2019', duration: 3, progress: 0.6 },
        { id: 2, text: 'Task #2', start_date: '18-04-2019', duration: 3, progress: 0.4 }
    ],
    links: [
        { id: 1, source: 1, target: 2, type: '0' }
    ]
};

class Gantt extends Component {
  componentDidMount() {
    gantt.init(this.ganttContainer);
    gantt.parse(data);
  }

  render() {
    return (
      <div
        ref={ (input) => { this.ganttContainer = input; } }
        style={ { width: '100%', height: '100%' } }
      ></div>
    );
  }
}

export default Gantt;
