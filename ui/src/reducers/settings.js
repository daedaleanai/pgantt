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

import {
  DATE_RANGE_SET, ZOOM_SET, SHOW_TASKS_OUTSIDE_TIMESCALE_SET,
  SHOW_TASKS_UNSCHEDULED_SET, SHOW_TASKS_CLOSED_SET
} from '../actions/settings';

// The initial UI settings.
const initialSettingsState = {
  startDate: null,
  endDate: null,
  zoom: "Days",
  showTasksOutsideTimescale: true,
  showTasksUnscheduled: false,
  showTasksClosed: false
};

// Performs the action onto the state and returns the updated state.
export function settingsReducer(state = initialSettingsState, action) {
  switch(action.type) {
  case DATE_RANGE_SET:
    return {
      ...state,
      startDate: action.startDate,
      endDate: action.endDate
    };

  case ZOOM_SET:
    return {
      ...state,
      zoom: action.zoom
    };

  case SHOW_TASKS_OUTSIDE_TIMESCALE_SET:
    return {
      ...state,
      showTasksOutsideTimescale: action.setting
    };

  case SHOW_TASKS_UNSCHEDULED_SET:
    return {
      ...state,
      showTasksUnscheduled: action.setting
    };

  case SHOW_TASKS_CLOSED_SET:
    return {
      ...state,
      showTasksClosed: action.setting
    };

  default:
    return state;
  }
}
