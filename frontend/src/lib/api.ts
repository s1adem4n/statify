const API_URL = import.meta.env.DEV ? 'http://localhost:8091' : window.location.origin;

type ViewsRequestOptions = {
	domain: string;
	start: Date;
	end?: Date;
	interval?: string;
};

export class Client {
	fetchViews(path: string, options: ViewsRequestOptions) {
		const url = `${API_URL}/stats/views/${path}`;
		// date as rfc3339
		const params = new URLSearchParams({
			domain: options.domain,
			start: options.start.toISOString()
		});
		if (options.end) {
			params.append('end', options.end.toISOString());
		}
		if (options.interval) {
			params.append('interval', options.interval);
		}

		return fetch(`${url}?${params.toString()}`).then((res) => res.json());
	}

	views = {
		count: async (options: ViewsRequestOptions) => {
			return (await this.fetchViews('count', options)) as number;
		},
		paths: async (options: ViewsRequestOptions) => {
			return (await this.fetchViews('paths', options)) as Record<string, number>;
		},
		devices: async (options: ViewsRequestOptions) => {
			return (await this.fetchViews('devices', options)) as Record<string, number>;
		},
		sessions: async (options: ViewsRequestOptions) => {
			return (await this.fetchViews('sessions', options)) as Record<string, number>;
		},
		time: async (options: ViewsRequestOptions) => {
			return (await this.fetchViews('time', options)) as number[];
		},
		domains: async () => {
			const url = `${API_URL}/stats/views/domains`;
			return (await fetch(url).then((res) => res.json())) as Record<string, number>;
		}
	};
}

const client = new Client();
export default client;
