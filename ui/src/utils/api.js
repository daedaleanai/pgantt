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

import { extractData } from './helpers';

const loc = window.location;
const api = process.env.NODE_ENV === 'production'
  ? `${loc.protocol}//${loc.host}/api`
  : 'http://localhost:9999/api';

const headers = {
  'Accept': 'application/json',
};

const responseHandler = (response) => {
  if(!response.ok) {
    if (response.error) {
      throw Error(response.error);
    }
    return response.json().then(body => {
      throw Error(body.data);
    });
  }
  return response.json();
};

export const projectsGet = () => {
  const url = `${api}/projects`;
  return fetch(url, { headers })
    .then(responseHandler);
};

export const planGet = (phid, closed) => {
  const url = `${api}/plan/${phid}?closed=${closed}`;
  return fetch(url, { headers })
    .then(responseHandler);
};

export const taskCreate = (phid, data) => {
  const url = `${api}/edit/${phid}/task`;
  return fetch(url, {
    method: "POST",
    headers,
    body: JSON.stringify(data)
  })
    .then(responseHandler)
    .then(extractData);
};

export const taskEdit = (phid, data) => {
  const url = `${api}/edit/${phid}/task`;
  return fetch(url, {
    method: "PUT",
    headers,
    body: JSON.stringify(data)
  })
    .then(responseHandler)
    .then(extractData);
};

export const taskDelete = (phid, id) => {
  return Promise.reject(new Error("You cannot delete tasks."));
};

export const linkCreate = (phid, data) => {
  const url = `${api}/edit/${phid}/link`;
  return fetch(url, {
    method: "POST",
    headers,
    body: JSON.stringify(data)
  })
    .then(responseHandler)
    .then(extractData);
};

export const linkEdit = (phid, data) => {
  return Promise.reject(new Error("You cannot edit links."));
};

export const linkDelete = (phid, id) => {
  const url = `${api}/edit/${phid}/link`;
  return fetch(url, {
    method: "DELETE",
    headers,
    body: JSON.stringify(id)
  })
    .then(responseHandler)
    .then(extractData);
};
