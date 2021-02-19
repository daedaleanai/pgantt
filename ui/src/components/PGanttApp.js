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

import React, { Component } from 'react';
import moment from 'moment';

import PGanttNav from './PGanttNav';
import ProjectView from './ProjectView';
import Gantt from './Gantt';
import GanttToolbar from './GanttToolbar';

class PGanttApp extends Component {
  render() {
    return (
      <div className="box">
        <div className="row header">
          <PGanttNav />
        </div>
        <ProjectView />
        <div className="row footer">
        </div>
      </div>
    );
  }
}

export default PGanttApp;
