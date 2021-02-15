import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import 'antd/dist/antd.css';
import 'dhtmlx-gantt/codebase/dhtmlxgantt.css';

import './index.css';
import PGanttApp from './components/PGanttApp';

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter>
      <PGanttApp />
    </BrowserRouter>
  </React.StrictMode>,
  document.getElementById('root')
);
