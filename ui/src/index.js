//------------------------------------------------------------------------------
// Copyright (C) 2021 Daedalean AG
//
// This file is part of PGantt.
//
// PGantt is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 2 of the License, or
// (at your option) any later version.
//
// PGantt is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with PGantt.  If not, see <https://www.gnu.org/licenses/>.
//------------------------------------------------------------------------------

import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { createStore, combineReducers } from 'redux';
import { Provider } from 'react-redux';

import 'antd/dist/antd.css';
import 'dhtmlx-gantt/codebase/dhtmlxgantt.css';

import './index.css';
import PGanttApp from './components/PGanttApp';

import { projectsReducer } from './reducers/projects';
import { planningReducer } from './reducers/planning';
import { settingsReducer } from './reducers/settings';

export const store = createStore(
  combineReducers({
    projects: projectsReducer,
    planning: planningReducer,
    settings: settingsReducer,
  }),
  window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__()
);

ReactDOM.render(
  <Provider store={store}>
    <BrowserRouter>
      <PGanttApp />
    </BrowserRouter>
  </Provider>,
  document.getElementById('root')
);
