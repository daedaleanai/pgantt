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

export const DATE_RANGE_SET = 'DATE_RANGE_SET';
export const ZOOM_SET = 'ZOOM_SET';
export const SHOW_TASKS_OUTSIDE_TIMESCALE_SET = 'SHOW_TASKS_OUTSIDE_TIMESCALE_SET';
export const SHOW_TASKS_UNSCHEDULED_SET = 'SHOW_TASKS_UNSCHEDULED_SET';
export const SHOW_TASKS_CLOSED_SET = 'SHOW_TASKS_CLOSED_SET';

// Creates an action that sets the date range setting.
export function dateRangeSet(startDate, endDate) {
  return {
    type: DATE_RANGE_SET,
    startDate,
    endDate
  };
}

// Creates an action that sets the zoom setting.
export function zoomSet(zoom) {
  return {
    type: ZOOM_SET,
    zoom
  };
}

// Creates an action that sets the showTasksOutsideTimescale setting.
export function showTasksOutsideTimescaleSet(setting) {
  return {
    type: SHOW_TASKS_OUTSIDE_TIMESCALE_SET,
    setting
  };
}

// Creates an action that sets the showTasksUnscheduled setting.
export function showTasksUnscheduledSet(setting) {
  return {
    type: SHOW_TASKS_UNSCHEDULED_SET,
    setting
  };
}

// Creates an action that sets the showTasksClosed setting.
export function showTasksClosedSet(setting) {
  return {
    type: SHOW_TASKS_CLOSED_SET,
    setting
  };
}
