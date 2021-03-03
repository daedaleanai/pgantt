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
import { gantt } from 'dhtmlx-gantt';
import { message } from 'antd';
import moment from 'moment';
import { connect } from 'react-redux';

import { planSet } from '../actions/planning';
import {
  planGet, taskCreate, taskEdit, taskDelete, linkCreate, linkEdit, linkDelete
} from '../utils/api';
import { objectEquals, sanitizeTask, sanitizeLink } from '../utils/helpers';

class Gantt extends Component {
  fetchData = (phid) => {
    planGet(phid, true)
      .then(data => this.props.planSet(data.data))
      .catch(msg => message.error(msg.toString()));
  }

  componentDidMount() {
    gantt.init(this.ganttContainer);
    gantt.plugins({
      marker: true
    });

    gantt.config.date_format = "%Y-%m-%d";

    const dateToStr = gantt.date.date_to_str(gantt.config.date_format);
    gantt.templates.format_date = (date) => {
      return dateToStr(date);
    };

    const strToDate = gantt.date.str_to_date(
      gantt.config.date_format, gantt.config.server_utc);

    gantt.templates.parse_date = (date) => {
      return strToDate(date);
    };

    var today = new Date();
    gantt.addMarker({
      start_date: today,
      css: "today",
      text: "Today",
      title: "Today: " + dateToStr(today)
    });

    gantt.i18n.setLocale({
      labels:{
        time_enable_button: "Schedule",
        time_disable_button: "Unschedule",
        section_details: "Details",
        section_title: "Title",
        section_column: "Column",
        section_parent: "Parent"
      }
    });

    const columns = this.props.project.columns.map((obj) => {
      return {key: obj.phid, label: obj.name};
    });

    const fields = [
      {name: "title", height: 70, map_to: "text", type: "textarea", focus: true},
      {name: "details", height: 16, type: "template", map_to: "details"},
      {name: "type", type: "typeselect", map_to: "type"},
      {name: "parent", type: "parent", allow_root: "true", root_label: "No parent"},
      {name: "column", height:22, map_to: "column", type: "select", options: columns},
      {name: "time", map_to: "auto", button: true, type: "duration_optional"}
    ];

    gantt.config.lightbox.sections = fields;
    gantt.config.lightbox.project_sections = fields;
    gantt.config.lightbox.milestone_sections = fields;

    gantt.templates.rightside_text = (start, end, task) => {
      if (task.type == gantt.config.types.milestone) {
        return task.text;
      }
      return "";
    };

    gantt.config.grid_width = 420;
    gantt.config.row_height = 24;
    gantt.config.grid_resize = true;

    gantt.config.columns = [
      {name: "text", tree: true, width: '*', resize: true},
      {name: "add", width: 40,  },
    ];

    gantt.attachEvent("onLightboxSave", (id, task, is_new) => {
      task.unscheduled = !task.start_date;
      return true;
    });

    gantt.attachEvent("onBeforeLightbox", (id) => {
      var task = gantt.getTask(id);
      task.details = `<b>URL:</b> <a href="${task.url}">${task.url}</a>`;
      return true;
    });

    gantt.config.auto_types = true;

    const logError = err => {
      message.error(err.message);
      throw err;
    };

    let dp = gantt.createDataProcessor({
      task: {
        create: (data) => {
          return taskCreate(this.props.phid, sanitizeTask(data))
            .catch(logError);
        },
        update: (data, id) => {
          return taskEdit(this.props.phid, sanitizeTask(data))
            .catch(logError);
        },
        delete: (id) => {
          return taskDelete(this.props.phid, id)
            .catch(logError);
        }
      },
      link: {
        create: (data) => {
          return linkCreate(this.props.phid, sanitizeLink(data))
            .catch(logError);
        },
        update: (data, id) => {
          return linkEdit(this.props.phid, sanitizeLink(data))
            .catch(logError);
        },
        delete: (id) => {
          return linkDelete(this.props.phid, id)
            .catch(logError);
        }
      }
    });

    dp.attachEvent("onAfterUpdate", (id, action, tid, response) => {
      if(action == "error") {
        gantt.clearAll();
        gantt.parse(this.props.plan);
      }
    });

    gantt.config.buttons_left = [];
    gantt.config.buttons_right = ["gantt_cancel_btn", "gantt_save_btn"];
    this.fetchData(this.props.phid);
    this.initGanttDataProcessor();
    setInterval(() => this.fetchData(this.props.phid), 1000);
  }

  componentWillReceiveProps(nextProps) {
    const thisPhid = this.props.phid;
    const nextPhid = nextProps.phid;
    if(thisPhid !== nextPhid) {
      this.fetchData(nextPhid);
    }
  }

  setZoom(value) {
    switch (value) {
    case 'Days':
      gantt.config.min_column_width = 70;
      gantt.config.scale_unit = 'week';
      gantt.config.date_scale = 'Week %W';
      gantt.config.subscales = [
        { unit: 'day', step: 1, date: '%d %M' }
      ];
      gantt.config.scale_height = 60;
      break;
    case 'Months':
      gantt.config.min_column_width = 70;
      gantt.config.scale_unit = 'month';
      gantt.config.date_scale = '%F';
      gantt.config.scale_height = 60;
      gantt.config.subscales = [
        { unit:'week', step:1, date:'Week %W' }
      ];
      break;
    default:
      break;
    }
  }

  setRange(start, end) {
    let s = start;
    let e = end;
    if (start == null || end === null) {
      s = moment();
      e = moment();
      if (s.day() !== 0) {
        s = s.day(0); // previous Sunday
      }
      if (e.day() !== 6) {
        e = e.day(6); // next Saturday
      }

    }
    gantt.config.start_date = s.toDate();
    gantt.config.end_date = e.toDate();
  }

  initGanttDataProcessor() {
    const onDataUpdated = this.props.onDataUpdated;
    this.dataProcessor = gantt.createDataProcessor((entityType, action, item, id) => {
      return new Promise((resolve, reject) => {
        if (onDataUpdated) {
         onDataUpdated(entityType, action, item, id);
        }
        return resolve();
      });
    });
  }

  componentWillUnmount() {
    if (this.dataProcessor) {
      this.dataProcessor.destructor();
      this.dataProcessor = null;
    }
  }

  shouldComponentUpdate(nextProps) {
    if (this.props.zoom !== nextProps.zoom) {
      return true;
    }

    if (this.props.startDate !== nextProps.startDate) {
      return true;
    }

    if (this.props.endDate !== nextProps.endDate) {
      return true;
    }

    if (this.props.plan.data.length !== nextProps.plan.data.length) {
      return true;
    }

    for (var i = 0; i < this.props.plan.data.length; i++) {
      if (!objectEquals(this.props.plan.data[i], nextProps.plan.data[i])) {
        return false;
      }
    }

    if (this.props.plan.links.length !== nextProps.plan.links.length) {
      return true;
    }

    for (i = 0; i < this.props.plan.links.length; i++) {
      if (!objectEquals(this.props.plan.links[i], nextProps.plan.links[i])) {
        return false;
      }
    }

    if (this.props.showTasksOutsideTimescale !== nextProps.showTasksOutsideTimescale) {
      return true;
    }

    return false;
  }

  componentDidUpdate() {
    gantt.config.show_tasks_outside_timescale = this.props.showTasksOutsideTimescale;
    gantt.parse(this.props.plan);
    gantt.render();
  }

  render() {
    const { zoom, startDate, endDate } = this.props;
    this.setRange(startDate, endDate);
    this.setZoom(zoom);
    return (
      <div
        ref={ (input) => { this.ganttContainer = input; } }
        style={ { width: '100%', height: '100%' } }
      ></div>
    );
  }
}

function mapStateToProps(state, ownProps) {
  const proj = state.projects.filter(proj => proj.phid === ownProps.phid);
  return {
    plan: state.planning,
    project: proj.length !== 0 ? proj[0] : null
  };
}

function mapDispatchToProps(dispatch) {
  return {
    planSet: (data) => dispatch(planSet(data))
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Gantt);
