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

const loc = window.location;
const api = process.env.NODE_ENV === 'production'
  ? `${loc.protocol}//${loc.host}/api`
  : 'http://localhost:9999/api';

const headers = {
  'Accept': 'application/json',
};

const responseHandler = (response) => {
  if(!response.ok)
    throw Error(response.error);
  return response.json();
};

export const projectsGet = () => {
  const url = `${api}/projects`;
  return fetch(url, { headers })
    .then(responseHandler);
};
