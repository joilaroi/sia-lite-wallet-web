<template>
	<div class="connect-ledger">
		<div class="step-icon"><icon :icon="['fab', 'usb']" /></div>
		<div class="title">{{ translate('ledger.instructions') }}</div>
		<button class="btn btn-inline btn-success" @click="onConnect" :disabled="connected">{{ translate('connect') }}</button>
	</div>
</template>

<script>
import { getVersion } from '@/ledger';

export default {
	props: {
		connected: Boolean
	},
	methods: {
		async onConnect() {
			try {
				await getVersion();
				this.$emit('connected', true);
			} catch (ex) {
				console.error('onConnect', ex);
				this.pushNotification({
					severity: 'danger',
					icon: ['fab', 'usb'],
					message: ex.message
				});
				this.$emit('connected', false);
			}
		}
	}
};
</script>

<style lang="stylus" scoped>
.connect-ledger {
	display: grid;
	grid-template-columns: auto minmax(0, 1fr) auto;
	grid-gap: 15px;
	padding: 15px;
	align-items: center;
	border-radius: 4px;
	background: dark-gray;

	svg {
		width: 32px;
		height: auto;
	}

	button {
		margin: 0;
	}
}
</style>