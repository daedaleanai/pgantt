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
  if (!response.ok) {
    if (response.error) {
      throw Error(response.error);
    }
    return response.json().then(body => {
      throw Error(body.data);
    });
  }
  return response.json();
};

// Get the list of projects from the webapp server.
export const projectsGet = () => {
  const url = `${api}/projects`;
  return fetch(url, { headers })
    .then(responseHandler);
};

// Get the tasks of the project from the webapp server.
export const planGet = (projectPhid) => {
  const url = `${api}/plan/${projectPhid}`;
  return fetch(url, { headers })
    .then(responseHandler);
};

// Tells the webapp server to create a task.
export const taskCreate = (projectPhid, task) => {
  const url = `${api}/edit/${projectPhid}/task`;
  return fetch(url, {
    method: "POST",
    headers,
    body: JSON.stringify(task)
  })
    .then(responseHandler)
    .then(extractData);
};

// Tells the webapp server to update a task.
export const taskEdit = (projectPhid, task) => {
  const url = `${api}/edit/${projectPhid}/task`;
  return fetch(url, {
    method: "PUT",
    headers,
    body: JSON.stringify(task)
  })
    .then(responseHandler)
    .then(extractData);
};

// Tells the webapp server to delete a task.
export const taskDelete = (projectPhid, id) => {
  return Promise.reject(new Error("You cannot delete tasks."));
};

// Tells the webapp server to create a link between two tasks.
export const linkCreate = (projectPhid, link) => {
  const url = `${api}/edit/${projectPhid}/link`;
  return fetch(url, {
    method: "POST",
    headers,
    body: JSON.stringify(link)
  })
    .then(responseHandler)
    .then(extractData);
};

// Tells the webapp server to update a link between two tasks.
export const linkEdit = (projectPhid, link) => {
  return Promise.reject(new Error("You cannot edit links."));
};

// Tells the webapp server to delete a link between two tasks.
export const linkDelete = (projectPhid, linkId) => {
  const url = `${api}/edit/${projectPhid}/link`;
  return fetch(url, {
    method: "DELETE",
    headers,
    body: JSON.stringify(linkId)
  })
    .then(responseHandler)
    .then(extractData);
};
