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
import { PageHeader, Radio } from 'antd';

class GanttToolbar extends Component {
  state = {
    currentZoom: 'Days',
  };

  handleZoomChange = (e) => {
    this.setState({
      currentZoom: e.target.value,
    });

    if (this.props.onZoomChange) {
      this.props.onZoomChange(e.target.value);
    }
  }

  render() {
    const options = [
      { label: 'Days', value: 'Days' },
      { label: 'Months', value: 'Months' },
    ];

    return (
      <div className="zoom-bar">
        <PageHeader
          ghost={false}
          title="Platforms"
          extra={[
            <Radio.Group
              options={options}
              onChange={this.handleZoomChange}
              value={this.state.currentZoom}
              optionType="button"
            />,
          ]}
        >
        </PageHeader>
      </div>
    );
  }
}

export default GanttToolbar;