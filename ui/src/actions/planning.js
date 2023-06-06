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

export const PLAN_SET = 'PLAN_SET';

// Returns an action for setting the planning of the current project in the store.
export function planSet(plan) {
  // taskId -> parentTaskId
  let parents = new Map();
  plan.data.forEach((task) => {
    parents.set(task.id, task.parent)
  });

  // taskId -> how many ancestor tasks it has
  let _levels = new Map();
  function level(taskId) {
    if (_levels.has(taskId)) {
      return _levels.get(taskId);
    }
    let l = 0;
    if (parents.has(taskId)) {
      l = 1 + level(parents.get(taskId));
    }
    _levels.set(taskId, l);
    return l;
  }

  function byLevelDateAgeFn(a, b) {
    // Order by number of ancestors.
    let levelDelta = level(a.id) - level(b.id);
    if (levelDelta != 0) {
      return levelDelta;
    }

    // Order by start date.
    // Replace "" with "Z" so the tasks with no starting date appear below.
    var aStartDate = a.start_date || "Z";
    var bStartDate = b.start_date || "Z";
    if (aStartDate < bStartDate) {
      return -1;
    } else if (aStartDate > bStartDate) {
      return 1;
    }
  
    // Order by task age.
    var aTaskIDNumber = parseInt(a.url.split("T")[1])
    var bTaskIDNumber = parseInt(b.url.split("T")[1])
    return aTaskIDNumber - bTaskIDNumber;
  }
  
  // The tasks order must be stable. Successive plans are compared to figure
  // out whether an expensive render is needed.
  // The sorting is done first by level (number of ancestors). Gantt will
  // iterate the sorted tasks to place the children under their parents, so
  // it all goes well in the end.
  plan.data.sort(byLevelDateAgeFn);

  return {
    type: PLAN_SET,
    plan
  };
}
