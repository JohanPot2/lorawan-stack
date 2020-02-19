// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React from 'react'
import { connect } from 'react-redux'
import { Col, Row, Container } from 'react-grid-system'
import bind from 'autobind-decorator'
import { defineMessages } from 'react-intl'

import PageTitle from '../../../components/page-title'
import LocationForm from '../../components/location-form'
import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'
import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import withFeatureRequirement from '../../lib/components/with-feature-requirement'
import sharedMessages from '../../../lib/shared-messages'

import { updateGateway } from '../../store/actions/gateways'
import attachPromise from '../../../lib/store/actions/attach-promise'
import { selectSelectedGateway, selectSelectedGatewayId } from '../../store/selectors/gateways'
import { mayViewOrEditGatewayLocation } from '../../lib/feature-checks'

import PropTypes from '../../../lib/prop-types'

const m = defineMessages({
  setGatewayLocation: 'Set gateway antenna location',
})

const getRegistryLocation = function(antennas) {
  let registryLocation
  if (antennas) {
    for (const key of Object.keys(antennas)) {
      if (antennas[key].location.source === 'SOURCE_REGISTRY') {
        registryLocation = { antenna: antennas[key], key }
        break
      }
    }
  }
  return registryLocation
}

@connect(
  state => ({
    gateway: selectSelectedGateway(state),
    gtwId: selectSelectedGatewayId(state),
  }),
  { updateGateway: attachPromise(updateGateway) },
)
@withBreadcrumb('gateway.single.data', function(props) {
  const { gtwId } = props
  return <Breadcrumb path={`/gateways/${gtwId}/location`} content={sharedMessages.location} />
})
@withFeatureRequirement(mayViewOrEditGatewayLocation, {
  redirect: ({ gtwId }) => `/gateways/${gtwId}`,
})
export default class GatewayLocation extends React.Component {
  static propTypes = {
    gateway: PropTypes.gateway.isRequired,
    gtwId: PropTypes.string.isRequired,
    updateGateway: PropTypes.func.isRequired,
  }

  @bind
  async handleSubmit(values) {
    const { gateway, gtwId, updateGateway } = this.props

    const patch = {}
    const registryLocation = getRegistryLocation(gateway.antennas)
    if (registryLocation) {
      // Update old location value
      patch.antennas = [...gateway.antennas]
      patch.antennas[registryLocation.key].location = {
        ...registryLocation.antenna.location,
        ...values,
      }
    } else {
      // Create new location value
      patch.antennas = [
        {
          gain: 0,
          location: {
            ...values,
            accuracy: 0,
            source: 'SOURCE_REGISTRY',
          },
        },
      ]
    }

    await updateGateway(gtwId, patch)
  }

  @bind
  async handleDelete() {
    const { gateway, gtwId, updateGateway } = this.props
    const registryLocation = getRegistryLocation(gateway.antennas)

    const patch = {
      antennas: [...gateway.antennas],
    }
    patch.antennas.splice(registryLocation.key, 1)

    await updateGateway(gtwId, patch)
  }

  render() {
    const { gateway, gtwId } = this.props
    const registryLocation = getRegistryLocation(gateway.antennas)
    const initialValues = registryLocation ? registryLocation.antenna.location : undefined

    return (
      <Container>
        <PageTitle title={sharedMessages.location} />
        <Row>
          <Col lg={8} md={12}>
            <LocationForm
              entityId={gtwId}
              formTitle={m.setGatewayLocation}
              initialValues={initialValues}
              onSubmit={this.handleSubmit}
              onDelete={this.handleDelete}
            />
          </Col>
        </Row>
      </Container>
    )
  }
}
