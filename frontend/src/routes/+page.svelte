<script lang="ts">

	import { BACKEND_API_URL } from '../constants/backend';

	let transactionFile: FileList | undefined | null = undefined;
	let toastWarning = {
		show: false,
		message: '',
		state: ''
	};

	const handleToast = (msg: string, state: string) => {
		toastWarning = {
			show: true,
			message: msg,
			state: state
		};

		setTimeout(() => {
			toastWarning = {
				show: false,
				message: '',
				state: ''
			};
		}, 3000);
	};

	const uploadTransactionFile = async () => {
		const formData = new FormData();
		if (!transactionFile || (transactionFile && transactionFile.length === 0)) {
			handleToast('Please select a file to upload', 'warning');
			return;
		}

		formData.append('file', transactionFile[0]);

		try {
			const res = await fetch(`${BACKEND_API_URL}/upload`, {
				method: 'POST',
				body: formData
			});

			if (res.ok) {
				handleToast('Successfully uploaded file', 'success');
			} else {
				handleToast('Failed to upload file', 'error');
			}
			// eslint-disable-next-line @typescript-eslint/no-unused-vars
		} catch (e) {
			handleToast('Failed to upload file', 'error');
		}
	};
</script>

<div class="lg:container mx-auto space-y-2">
	<h1>Welcome to Funances!</h1>

	<label class="form-control w-full max-w-xs">
		<div class="label">
			<span class="label-text">Upload Transactions</span>
		</div>
		<input
			accept="text/plain, text/csv"
			bind:files={transactionFile}
			class="file-input file-input-bordered file-input-primary w-full max-w-xs"
			id="transactionFileInput"
			type="file" />
	</label>
	<button class="btn btn-primary" on:click={() => uploadTransactionFile()}>Upload</button>
</div>

{#if toastWarning.show}
	<div class="toast toast-top toast-start">
		{#if toastWarning.state === 'success'}
			<div class="alert alert-success">
				<span>{toastWarning.message}</span>
			</div>
		{:else if toastWarning.state === 'error'}
			<div class="alert alert-error">
				<span>{toastWarning.message}</span>
			</div>
		{:else if toastWarning.state === 'warning'}
			<div class="alert alert-warning">
				<span>{toastWarning.message}</span>
			</div>
		{/if}
	</div>
{/if}



