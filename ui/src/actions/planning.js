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

function byDateAgeFn(a, b) {
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

export function planSet(plan) {
  plan.data.sort(byDateAgeFn);

  return {
    type: PLAN_SET,
    plan
  };
}
